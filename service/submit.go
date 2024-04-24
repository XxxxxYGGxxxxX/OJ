package service

import (
	"OJ/define"
	"OJ/helper"
	"OJ/models"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetSubmitList
// @Tags 公共方法
// @Summary 提交列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param problem_identity query string false "problem_identity"
// @Param user_identity query string false "user_identity"
// @Param status query int false "status"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /submit-list [get]
func GetSubmitList(c *gin.Context) {
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("GetProblemList Page strconv Error:", err)
		return
	}
	page = (page - 1) * size
	var count int64
	list := make([]models.SubmitBasic, 0)
	problemIdentity := c.Query("problem_identity")
	userIdentity := c.Query("user_identity")
	status, _ := strconv.Atoi(c.Query("status"))
	tx := models.GetSubmitList(problemIdentity, userIdentity, status)
	err = tx.Count(&count).Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		log.Println("Get Problem List Error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Submit List Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}

// Submit
// @Tags 用户私有方法
// @Summary 代码提交
// @Param authorization header string true "authorization"
// @Param problem_identity query string true "problem_identity"
// @Param code body string true "code"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user/submit [post]
func Submit(c *gin.Context) {
	// 获取信息
	problemIdentity := c.Query("problem_identity")
	code, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Read Code Error:" + err.Error(),
		})
		return
	}
	// 代码保存
	path, err := helper.CodeSave(code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Code Save Error:" + err.Error(),
		})
		return
	}
	// 取用户标识
	u, _ := c.Get("user")
	userClaim := u.(*helper.UserClaims)
	sb := &models.SubmitBasic{
		Identity:        helper.GetUUID(),
		ProblemIdentity: problemIdentity,
		UserIdentity:    userClaim.Identity,
		Path:            path,
	}
	// 代码判断
	pb := new(models.ProblemBasic)
	// 答案错误的channel
	WA := make(chan int)
	// 超内存的channel
	OOM := make(chan int)
	// 编译错误的channel
	CE := make(chan int)
	// 通过次数
	passCount := 0
	// 互斥锁 用于passCount
	var lock sync.Mutex
	// 提示信息
	var msg string
	err = models.DB.Where("identity = ?", problemIdentity).Preload("TestCases").First(pb).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Problem Error:" + err.Error(),
		})
		return
	}
	for _, testCase := range pb.TestCase {
		testCase := testCase
		go func() {
			cmd := exec.Command("go", "run", path)
			// 标准输出、标准错误
			var out, stderr bytes.Buffer
			cmd.Stderr = &stderr
			cmd.Stdout = &out
			stdinPipe, err := cmd.StdinPipe()
			if err != nil {
				log.Fatalln(err)
			}
			// 写入输入案列
			io.WriteString(stdinPipe, testCase.Input)
			// MemStats是一个结构体,对用户提供只读的内存统计信息
			// 运行之前的内存
			var bm runtime.MemStats
			runtime.ReadMemStats(&bm)
			// 开始执行
			// 判断是否编译错误
			if err := cmd.Run(); err != nil {
				log.Println(err, stderr.String())
				if err.Error() == "exit status 2" {
					msg = stderr.String()
					CE <- 1
					return
				}
			}
			// 运行之后的内存
			var em runtime.MemStats
			runtime.ReadMemStats(&em)

			// 判断答案错误
			if testCase.Output != out.String() {
				msg = "答案错误"
				WA <- 1
				return
			}
			// 判断超内存
			if em.Alloc/1024-(bm.Alloc/1024) > uint64(pb.MaxMem) {
				msg = "运行超内存"
				OOM <- 1
				return
			}
			lock.Lock()
			passCount++
			lock.Unlock()
			// println("Err:", string(stderr.Bytes()))
			// fmt.Println(out.String())

			// println(out.String() == "34\n")
		}()
	}
	// 阻塞判断
	// status -1待判断 1正确 2错误 3超时 4超内存 5编译错误
	select {
	case <-WA:
		sb.Status = 2
	case <-OOM:
		sb.Status = 4
	case <-CE:
		sb.Status = 5
		// 判断是否运行时过长
	case <-time.After(time.Millisecond * time.Duration(pb.MaxRuntime)):
		// 通过的次数=测试案列的个数
		if passCount == len(pb.TestCase) {
			sb.Status = 1
		} else {
			sb.Status = 3
		}
	}
	// 开启事务,因为涉及多个表的修改
	if err = models.DB.Transaction(func(tx *gorm.DB) error {
		// 用户提供的代码插入到submit_basic表中
		err = tx.Create(sb).Error
		if err != nil {
			return errors.New("SubmitBasic Save Error:" + err.Error())
		}
		m := make(map[string]interface{})
		// 提交次数加一
		m["submit_num"] = gorm.Expr("submit_num + ?", 1)
		// 如果答案正确则通过次数加一
		if sb.Status == 1 {
			m["pass_num"] = gorm.Expr("pass_num + ?", 1)
		}
		// 更新user_basic表
		err = tx.Model(new(models.UserBasic)).Where("identity = ?", userClaim.Identity).Updates(m).Error
		if err != nil {
			return errors.New("SubmitBasic Modify Error:" + err.Error())
		}
		// 更新problem_basic表
		err = tx.Model(new(models.ProblemBasic)).Where("problem_identity = ?", problemIdentity).Updates(m).Error
		if err != nil {
			return errors.New("ProblemBasic Modify Error:" + err.Error())
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Submit Error:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"status": sb.Status,
			"msg":    msg,
		},
	})
}
