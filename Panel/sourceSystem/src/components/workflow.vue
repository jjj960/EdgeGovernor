<template>
  <div>
    <!-- style="position: absolute; top:0%" 很重要，可以解决下拉框偏移问题-->
    <el-dialog
      title="Workflow Deployment"
      :visible.sync="detailVisible"
      width="40%"
      style="position: absolute; top: -5%"
    >
      <el-row>
        <!-- 表单 -->
        <el-form :model="formInline" ref="form" label-width="160px">
          <el-form-item label="Workflow Name:" class="choose itemCenter">
            <el-input
              v-model="taskName"
              placeholder="Enter workflow name"
              class="input"
            ></el-input>
          </el-form-item>
          <el-form-item label="Job number:" class="choose itemCenter">
            <el-input
              v-model="taskNum"
              placeholder="Enter job number"
              class="input"
            ></el-input>
          </el-form-item>
          <el-form-item label="Workflow Deployment Mode：" class="choose itemCenter">
            <el-select
              v-model="tpValue"
              placeholder="Select deployment time"
              class="doubleInput"
              style="margin-bottom: 10px"
              v-on:change="choosepostTime(tpValue)"
            >
              <el-option
                v-for="item in taskpostValue"
                :key="item.tpValue"
                :label="item.tpLabel"
                :value="item.tpValue"
              >
              </el-option>
            </el-select>
            <el-input
              v-model="taskpostTime"
              class="doubleInput"
              :disabled="true"
            ></el-input>
          </el-form-item>
          <el-form-item label="Upload Method:" class="choose itemCenter">
            <el-select
              v-model="pmValue"
              placeholder="Select upload method"
              class="doubleInput"
              style="margin-bottom: 10px"
              v-on:change="choosepostMethod(pmValue)"
            >
              <el-option
                v-for="item in postMethods"
                :key="item.pmValue"
                :label="item.pmLabel"
                :value="item.pmValue"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="Text Upload：" class="choose itemCenter">
            <el-input
              type="textarea"
              placeholder="Enter context"
              v-model="textarea"
              :disabled="textPost"
            >
            </el-input>
          </el-form-item>
          <el-form-item label="File Upload:" class="choose itemCenter">
            <el-upload
              class="upload-demo"
              drag
              action="''"
              :limit="1"
              :on-remove="removeFile"
              :on-change="postFile"
              :disabled="filePost"
              :file-list="fileList"
              :auto-upload="false"
            >
              <i class="el-icon-upload"></i>
              <div class="el-upload__text">Drag files here,or<em>Click to upload</em></div>
              <div class="el-upload__tip" slot="tip" style="margin-top: -10px">
                Note: Up to one file can be uploaded simultaneously
              </div>
            </el-upload>
          </el-form-item>
        </el-form>
      </el-row>
      <el-row style="text-align: right; background-color: #253a4f">
        <el-button @click="disappear()" style="margin-left: -80px">Cancel</el-button>
        <el-button type="primary" @click="workLoadPublish()">Confirm</el-button>
      </el-row>
    </el-dialog>
    <!-- 窗口 -->
    <taskpost ref="msgBtn" v-if="Visiable" @select="selectTime"></taskpost>
  </div>
</template>

<script>
import taskpost from "./taskpost";
import filter from "./filter.js";
import axios from "axios";
axios.defaults.baseURL = "/api";
export default {
  data() {
    return {
      detailVisible: false,

      // 与弹出窗口显示有关
      Visiable: false,

      taskName: "",
      taskNum: "",
      // 任务发布模式下拉框
      taskpostValue: [
        {
          tpValue: "scheRelease",
          tpLabel: "Scheduled Deployment",
        },
        {
          tpValue: "insRelease",
          tpLabel: "Instant Deployment",
        },
      ],
      tpValue: "",
      // 数据上传方式下拉框
      postMethods: [
        {
          pmValue: "file",
          pmLabel: "File Upload",
        },
        {
          pmValue: "text",
          pmLabel: "Text Upload",
        },
      ],
      pmValue: "",
      // 确保对应的控件能够使用
      filePost: true,
      textPost: true,
      // 任务发布时间
      date1: "",
      date2: "",
      taskpostTime: "",

      formInline: {},
      // 延时
      timer: null,

      formData: "",

      textarea: "",

      fileList: [],
    };
  },
  methods: {
    init(data) {
      this.detailVisible = true;
      //data是父组件弹窗传递过来的值，我们可以打印看看
      console.log(data);
    },
    disappear() {
      this.fileList = [];
      this.detailVisible = false;
    },
    // 获取会话内容
    getFromSessionStorage(key) {
      return window.sessionStorage.getItem(key);
    },
    // 选择发布时间
    choosepostTime(tpValue) {
      if (tpValue == "scheRelease") {
        this.examineBtn();
      } else if (tpValue == "insRelease") {
        this.taskpostTime = "";
      }
    },
    // 选择上传方式
    choosepostMethod(pmValue) {
      console.log(pmValue);
      if (pmValue == "file") {
        this.filePost = false;
        this.textPost = true;
      } else {
        this.filePost = true;
        this.textPost = false;
      }
    },
    examineBtn(data) {
      this.Visiable = true;
      this.$nextTick(() => {
        this.$refs.msgBtn.init(data);
      });
    },
    selectTime(date1, date2) {
      this.date1 = filter.formatDate(date1, "yyyy-MM-dd");
      this.date2 = filter.formatDate(date2, "HH:mm");
      this.taskpostTime = this.date1 + " " + this.date2;
    },
    // 移除文件
    removeFile(file) {
      console.log(file);
      const that = this;
      let ipAddress = that.getFromSessionStorage('url')
      let nodeUrl = `http://${ipAddress}`
      axios
        .post(
          `${nodeUrl}:5000/removeFile`,
          that.$qs.stringify({
            fileName: file.name,
          })
        )
        .then(function (response) {
          console.log(response);
          that.fileList = [];
        })
        .catch(function (error) {
          console.log(error);
        });
    },
    postFile(file) {
      const that = this;
      const formData = new FormData();
      formData.append("file", file.raw);
      let ipAddress = that.getFromSessionStorage('url')
      let nodeUrl = `http://${ipAddress}`
      // 这里的URL应该是你的文件上传接口
      axios
        .post(`${nodeUrl}:5000/upload`, formData)
        .then((response) => {
          // 处理成功响应
          console.log(response);
          that.fileList.push({ name: file.name, url: file.url });
        })
        .catch((error) => {
          // 处理错误
          console.error(error);
        });
    },
    // 提交任务
    workLoadPublish() {
      const that = this;
      let ipAddress = that.getFromSessionStorage('url')
      let nodeUrl = `http://${ipAddress}`
      axios
        .post(
          `${nodeUrl}:5000/workLoadPublish`,
          that.$qs.stringify({
            taskName: that.taskName,
            taskNum: that.taskNum,
            // 任务发布模式，定时发的时间(taskpostTime1代表日期，taskpostTime2代表小时分钟)
            tpValue: that.tpValue,
            taskpostTime1: that.date1,
            taskpostTime2: that.date2,
            pmValue: that.pmValue,
            textArea: that.textArea,
            fileList: JSON.stringify(that.fileList),
            date: filter.formatDate(new Date(), "yyyy-MM-dd HH:mm:ss"),
          })
        )
        .then(function (response) {
          that.$message({
            showClose: true,
            type: "success",
            message: "Publish Success!",
          });
          console.log(response);
          that.disappear();
        })
        .catch(function (error) {
          that.$message({
            showClose: true,
            type: "warning",
            message: "Publish Fail!",
          });
          console.log(error);
        });
    },
  },
  components: {
    taskpost,
  },
};
</script>

<style scoped>
.choose {
  margin-top: 10px;
  font-size: 14px;
  font-weight: bolder;
  color: #555555;
  width: 100%;
}

.input {
  width: 100%;
}

.doubleInput {
  width: 100%;
}

@media screen and (max-width: 1900px) {
  /deep/ .el-form-item__label {
    font-size: 15px;
  }

  /deep/ .el-input {
    position: relative;
    font-size: 14px;
  }

  /deep/ .el-form-item__content {
    width: 66%;
  }
}
</style>
