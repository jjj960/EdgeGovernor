<!-- 选择任务发布时间的窗口 -->
<template>
  <!--1.首先，弹窗页面中要有el-dialog组件即弹窗组件，我们把弹窗中的内容放在el-dialog组件中-->

  <!--2.设置:visible.sync属性，动态绑定一个布尔值，通过这个属性来控制弹窗是否弹出-->
  <el-dialog
    title="Alarm Message"
    :visible.sync="detailVisible"
    width="40%"
    style="top: 5%; position: absolute"
  >
    <div class="main_table">
      <!-- 查询 -->
      <div>
        <el-form :inline="true" :model="searchArea2">
          <el-form-item style="margin-left: 4px; margin-bottom: 3px">
            <el-button type="primary" @click="toggleSelection()">Select Cancel</el-button>
          </el-form-item>
          <el-form-item style="margin-left: 4px; margin-bottom: 3px">
            <el-button type="danger" @click="deleteWindow('NotAll')">Batch Delete</el-button>
          </el-form-item>
          <el-form-item style="margin-left: 4px; margin-bottom: 3px">
            <el-button type="warning" @click="deleteWindow('All')">Delete All</el-button>
          </el-form-item>
          <el-form-item style="margin-left: 4px; margin-bottom: 3px">
            <el-button type="success" @click="refresh()">Refresh</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表格 -->

      <div class="box-content">
        <el-table
          ref="multipleTable"
          :data="tableData.slice((currentPage - 1) * pageSize, currentPage * pageSize)"
          style="width: 100%"
          @selection-change="selectionChange"
        >
          <el-table-column type="selection" width="55"> </el-table-column>
          <el-table-column prop="id" label="ID" width="80"> </el-table-column>
          <el-table-column prop="type" label="Alarm Type" width="180"> </el-table-column>
          <el-table-column prop="date" label="Alarm Time" width="180"> </el-table-column>
          <el-table-column prop="detail" label="Detail"> </el-table-column>
        </el-table>
      </div>
      <!-- 页码 -->
      <div class="block">
        <el-pagination
          layout="total, prev, pager, next"
          :total="pageSizes"
          :page-size="pageSize"
          @current-change="handleCurrentChange"
          :current-page.sync="currentPage"
        >
        </el-pagination>
      </div>
    </div>
  </el-dialog>
</template>

<script>
import axios from "axios";
export default {
  name: "warningMes",
  data() {
    return {
      detailVisible: false,
      searchArea2: {},
      pageSizes: 0,
      pageSize: 4,
      currentPage: 1,
      tableData: [],
      selectionList: [],
      timer: null,
    };
  },
  methods: {
    //3.定义一个init函数，通过设置detailVisible值为true来让弹窗弹出，这个函数会在父组件的方法中被调用
    init(data) {
      this.detailVisible = true;
      //data是父组件弹窗传递过来的值，我们可以打印看看
      this.getTaskLog();
    },
    disappear() {
      this.detailVisible = false;
    },
    // 切换页码
    handleCurrentChange(val) {
      this.currentPage = val;
      // this.getnodeMessage1()
      // this.generateRelate()
    },
    // 获取会话内容
    getFromSessionStorage(key) {
      return window.sessionStorage.getItem(key);
    },
    // 请求节点信息
    getTaskLog() {
      const that = this;
      let ipAddress = that.getFromSessionStorage('url')
      let nodeUrl = `http://${ipAddress}`
      axios
        .get(`${nodeUrl}:5000/getTaskLog`, {
          data: "",
        })
        .then(function (response) {
          console.log(response.data);
          that.tableData = response.data.taskLogs;
          that.pageSizes = response.data.pageSizes;
        })
        .catch(function (error) {
          console.log(error);
        });
    },
    refresh() {
      this.currentPage = 1;
      this.getTaskLog();
    },
    selectionChange(val) {
      this.selectionList = [];
      val.forEach((element) => {
        this.selectionList.push(element.id);
      });
    },
    toggleSelection(rows) {
      if (rows) {
        rows.forEach((row) => {
          this.$refs.multipleTable.toggleRowSelection(row);
        });
      } else {
        this.$refs.multipleTable.clearSelection();
      }
    },
    // 删除任务确定窗口
    deleteWindow(data) {
      this.$confirm("Confirm deletion", "Tips", {
        confirmButtonText: "Confirm",
        cancelButtonText: "Cancel",
        type: "warning",
      })
        .then(() => {
          if (data !== "All") {
            this.handlerDelete();
            clearTimeout(this.timer); //清除延迟执行
            this.timer = setTimeout(() => {
              //设置延迟执行
              this.refresh();
            }, 500);
          } else {
            this.deleteAllLog();
            clearTimeout(this.timer); //清除延迟执行
            this.timer = setTimeout(() => {
              //设置延迟执行
              this.refresh();
            }, 500);
          }
        })
        .catch((error) => {
          console.log(error);
        });
    },
    handlerDelete() {
      const that = this;
      // 数组转字符串使用逗号分隔
      if (this.selectionList.length === 0) {
        this.$message({
          showClose: true,
          type: "warning",
          message: "No information selected",
        });
        return;
      }
      let sids = this.selectionList.join(",");
      let ipAddress = that.getFromSessionStorage('url')
      let nodeUrl = `http://${ipAddress}`
      axios
        .post(
          `${nodeUrl}:5000/deleteTaskLog`,
          that.$qs.stringify({
            sids: sids,
          })
        )
        .then(function (response) {
          that.$message({
            showClose: true,
            type: "success",
            message: "Delete Success!",
          });
        })
        .catch(function (error) {
          console.log(error);
          that.$message({
            showClose: true,
            type: "danger",
            message: "Delete Fail!",
          });
        });
    },
    deleteAllLog() {
      const that = this;
      let ipAddress = that.getFromSessionStorage('url')
      let nodeUrl = `http://${ipAddress}`
      // 数组转字符串使用逗号分隔
      axios
        .get(`${nodeUrl}:5000/deleteAllLog`, {
          data: "",
        })
        .then(function (response) {
          that.$message({
            showClose: true,
            type: "success",
            message: "Delete Success!",
          });
          that.refresh();
        })
        .catch(function (error) {
          console.log(error);
        });
    },
  },
};
</script>
<style scoped>
.el-form-item__content {
  width: 100% !important;
}

.main_table {
  text-align: left;
  margin-top: -15px;
}

.block {
  text-align: center;
}
</style>
