<template>
  <div>
    <div class="header">
      <div class="bg_header">
        <div class="header_nav fl t_title" style="margin-left: 3%">
          EdgeGovernor
          <span style="font-size: 12px">Current Leader:{{ url }}</span>
        </div>
      </div>
      <div class="message">
        <el-badge :value="notReadMes" class="item" type="warning">
          <div @click="examineBtn4()" style="cursor: pointer">
            <i class="el-icon-warning-outline"></i>
          </div>
        </el-badge>
      </div>
    </div>

    <!--main-->
    <div class="data_content">
      <div class="data_main">
        <div class="main_left fl">
          <div class="left_1 t_btn6" style="cursor: pointer">
            <!--左上边框-->
            <div class="t_line_box">
              <i class="t_l_line"></i>
              <i class="l_t_line"></i>
            </div>
            <!--右上边框-->
            <div class="t_line_box">
              <i class="t_r_line"></i>
              <i class="r_t_line"></i>
            </div>
            <!--左下边框-->
            <div class="t_line_box">
              <i class="l_b_line"></i>
              <i class="b_l_line"></i>
            </div>
            <!--右下边框-->
            <div class="t_line_box">
              <i class="r_b_line"></i>
              <i class="b_r_line"></i>
            </div>
            <div class="main_title">
              Task Distribution
              <span @click="gettaskNum()"><i class="el-icon-refresh"></i></span>
            </div>
            <div
              id="chart_1"
              class="chart"
              style="width: 100%; height: 250px; margin-top: 5px"
            ></div>
          </div>
          <div
            class="left_2"
            style="
              width: 100%;
              height: 100%;
              box-sizing: border-box;
              border: 1px solid #2c58a6;
              position: relative;
              box-shadow: 0 0 10px #2c58a6;
            "
          >
            <!--左上边框-->
            <div class="t_line_box">
              <i class="t_l_line"></i>
              <i class="l_t_line"></i>
            </div>
            <!--右上边框-->
            <div class="t_line_box">
              <i class="t_r_line"></i>
              <i class="r_t_line"></i>
            </div>
            <!--左下边框-->
            <div class="t_line_box">
              <i class="l_b_line"></i>
              <i class="b_l_line"></i>
            </div>
            <!--右下边框-->
            <div class="t_line_box">
              <i class="r_b_line"></i>
              <i class="b_r_line"></i>
            </div>
            <div class="main_title" style="cursor: pointer">
              Resource Usage
              <span @click="sourceRefresh()"><i class="el-icon-refresh"></i></span>
            </div>
            <div
              class="core_dialogue"
              style="
                width: 100%;
                height: 330px;
                margin-top: 30px;
                border-bottom: 1px solid rgb(44, 88, 166);
                overflow-y: scroll;
                overflow-x: hidden;
              "
            >
              <el-row
                style="
                  font-size: 18px;
                  color: white;
                  font-weight: 600;
                  margin-bottom: 15px;
                "
                >CPU Usage</el-row
              >
              <el-row
                v-for="item in nodecpuSources"
                :key="item.nodeName"
                style="margin-right: 10px; text-align: center"
              >
                <div
                  style="
                    margin-bottom: 5px;
                    color: white;
                    font-weight: 600;
                    text-align: left;
                    margin-left: 7px;
                  "
                >
                  {{ item.nodeName }}:
                </div>
                <el-progress
                  :percentage="item.percentage"
                  :color="customColorMethod"
                  :stroke-width="11"
                ></el-progress>
              </el-row>
            </div>
            <div
              class="core_dialogue"
              style="
                width: 100%;
                height: 330px;
                margin-top: 30px;
                border-bottom: 1px solid rgb(44, 88, 166);
                overflow-y: scroll;
                overflow-x: hidden;
              "
            >
              <el-row
                style="
                  font-size: 18px;
                  color: white;
                  font-weight: 600;
                  margin-bottom: 15px;
                "
                >Mem Usage</el-row
              >
              <el-row
                v-for="item in nodememorySources"
                :key="item.nodeName"
                style="margin-right: 10px; text-align: center"
              >
                <div
                  style="
                    margin-bottom: 5px;
                    color: white;
                    font-weight: 600;
                    text-align: left;
                    margin-left: 7px;
                  "
                >
                  {{ item.nodeName }}:
                </div>
                <el-progress
                  :percentage="item.percentage"
                  :color="customColorMethod"
                  :stroke-width="11"
                ></el-progress>
              </el-row>
            </div>
            <div
              class="core_dialogue"
              style="
                width: 100%;
                height: 330px;
                margin-top: 30px;
                overflow-y: scroll;
                overflow-x: hidden;
              "
            >
              <el-row
                style="
                  font-size: 18px;
                  color: white;
                  font-weight: 600;
                  margin-bottom: 15px;
                "
                >Disk Usage</el-row
              >
              <el-row
                v-for="item in nodediskSources"
                :key="item.nodeName"
                style="margin-right: 10px; text-align: center"
              >
                <div
                  style="
                    margin-bottom: 5px;
                    color: white;
                    font-weight: 600;
                    text-align: left;
                    margin-left: 7px;
                  "
                >
                  {{ item.nodeName }}:
                </div>
                <el-progress
                  :percentage="item.percentage"
                  :color="customColorMethod"
                  :stroke-width="11"
                ></el-progress>
              </el-row>
            </div>
          </div>
        </div>
        <div class="main_center fl">
          <div class="center_text" style="height: 400px">
            <!--左上边框-->
            <div class="t_line_box">
              <i class="t_l_line"></i>
              <i class="l_t_line"></i>
            </div>
            <!--右上边框-->
            <div class="t_line_box">
              <i class="t_r_line"></i>
              <i class="r_t_line"></i>
            </div>
            <!--左下边框-->
            <div class="t_line_box">
              <i class="l_b_line"></i>
              <i class="b_l_line"></i>
            </div>
            <!--右下边框-->
            <div class="t_line_box">
              <i class="r_b_line"></i>
              <i class="b_r_line"></i>
            </div>
            <div class="main_title">Node Management</div>
            <!-- <div id="chart_map" style="width:100%;height:610px;"></div> -->
            <div class="bottom_center fl">
              <div class="main_table t_btn8">
                <!-- 查询 -->
                <div class="search">
                  <el-form :inline="true" :model="searchArea1" class="choose">
                    <el-form-item label="Node Search:" style="margin-bottom: 5px">
                      <el-select
                        v-model="nodeValue1"
                        placeholder="Select Node"
                        class="input"
                        @focus="getnodeName1()"
                      >
                        <el-option
                          v-for="item in nodes1"
                          :key="item.nodeValue1"
                          :label="item.nodeLabel1"
                          :value="item.nodeValue1"
                        >
                        </el-option>
                      </el-select>
                    </el-form-item>
                    <el-form-item style="margin-left: 20px; margin-bottom: 5px">
                      <el-button type="primary" @click="searchNode()">Search</el-button>
                    </el-form-item>
                    <el-form-item style="margin-left: 10px; margin-bottom: 5px">
                      <el-button type="success" @click="refreshNode()">Refresh</el-button>
                    </el-form-item>
                  </el-form>
                </div>
                <!-- 表格 -->
                <div class="box-content">
                  <el-table :data="tableData1" style="width: 100%">
                    <el-table-column prop="nodeName" label="Name" style="width: 12%">
                    </el-table-column>
                    <el-table-column prop="nodeType" label="Type" style="width: 12%">
                    </el-table-column>
                    <el-table-column
                      prop="cpuSize"
                      label="CPU Usage(m)"
                      style="width: 12%"
                    >
                    </el-table-column>
                    <el-table-column
                      prop="memorySize"
                      label="Mem Usage(Mb)"
                      style="width: 12%"
                    >
                    </el-table-column>
                    <el-table-column
                      prop="diskSize"
                      label="Disk Usage(Mb)"
                      style="width: 12%"
                    >
                    </el-table-column>
                    <el-table-column prop="IP" label="IP" style="width: 12%">
                    </el-table-column>
                    <el-table-column prop="status" label="Status" style="width: 12%">
                    </el-table-column>
                  </el-table>
                </div>
                <!-- 页码 -->
                <div class="block">
                  <el-pagination
                    layout="total, prev, pager, next"
                    :total="pageSizes1"
                    :page-size="pageSize"
                    @current-change="handleCurrentChange"
                    :current-page.sync="currentPage1"
                  >
                  </el-pagination>
                </div>
              </div>
            </div>
          </div>
          <div class="center_text" style="margin-top: 40px; height: 440px">
            <!--左上边框-->
            <div class="t_line_box">
              <i class="t_l_line"></i>
              <i class="l_t_line"></i>
            </div>
            <!--右上边框-->
            <div class="t_line_box">
              <i class="t_r_line"></i>
              <i class="r_t_line"></i>
            </div>
            <!--左下边框-->
            <div class="t_line_box">
              <i class="l_b_line"></i>
              <i class="b_l_line"></i>
            </div>
            <!--右下边框-->
            <div class="t_line_box">
              <i class="r_b_line"></i>
              <i class="b_r_line"></i>
            </div>
            <div class="main_title">Task Management</div>
            <div class="bottom_center fl">
              <div class="main_table t_btn8">
                <!-- 查询 -->
                <div class="search">
                  <el-form
                    :inline="true"
                    :model="searchArea2"
                    class="demo-form-inline choose"
                  >
                    <el-form-item label="Task Search:" style="margin-bottom: 5px">
                      <el-select
                        v-model="nodeValue2"
                        placeholder="Select Task"
                        class="input"
                        @focus="getnodeName2()"
                      >
                        <el-option
                          v-for="item in nodes2"
                          :key="item.nodeValue2"
                          :label="item.nodeLabel2"
                          :value="item.nodeValue2"
                        >
                        </el-option>
                      </el-select>
                    </el-form-item>
                    <el-form-item style="margin-left: 20px; margin-bottom: 5px">
                      <el-button type="primary" @click="searchTask()">Search</el-button>
                    </el-form-item>
                    <el-form-item style="margin-left: 10px; margin-bottom: 5px">
                      <el-button type="warning" @click="examineBtn1()"
                        >Task Deploy</el-button
                      >
                    </el-form-item>
                    <el-form-item style="margin-left: 10px; margin-bottom: 5px">
                      <el-button type="success" @click="examineBtn6()"
                        >Workflow Deploy</el-button
                      >
                    </el-form-item>
                    <el-form-item style="margin-bottom: 5px; margin-left: 20px">
                      <el-tag type="info">Current Node：{{ this.currentNode }}</el-tag>
                    </el-form-item>
                  </el-form>
                </div>
                <!-- 表格 -->
                <div class="box-content">
                  <el-table :data="tableData2" style="width: 100%">
                    <el-table-column prop="taskID" label="ID" style="width: 11%">
                    </el-table-column>
                    <el-table-column prop="taskName" label="Task Name" style="width: 11%">
                    </el-table-column>
                    <el-table-column
                      prop="taskMirror"
                      label="Task Mirror"
                      style="width: 11%"
                    >
                    </el-table-column>
                    <el-table-column
                      prop="cpuSize"
                      label="CPU Usage(m)"
                      style="width: 11%"
                    >
                    </el-table-column>
                    <el-table-column
                      prop="memorySize"
                      label="Mem Usage(Mb)"
                      style="width: 11%"
                    >
                    </el-table-column>
                    <el-table-column
                      prop="diskSize"
                      label="Disk Usage(Mb)"
                      style="width: 11%"
                    >
                    </el-table-column>
                    <el-table-column prop="network" label="Type" style="width: 11%">
                    </el-table-column>
                    <el-table-column label="Operation" style="width: 11%">
                      <template slot-scope="scope">
                        <el-button @click="examineBtn(scope.row)" type="text" size="small"
                          >Migrate</el-button
                        >
                        <el-button
                          @click="deleteWindow(scope.row)"
                          type="text"
                          size="small"
                          >Delete</el-button
                        >
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
                <!-- 页码 -->
                <div class="block">
                  <el-pagination
                    layout="total, prev, pager, next"
                    :total="pageSizes2"
                    :page-size="pageSize"
                    @current-change="searchTask()"
                    :current-page.sync="currentPage2"
                  >
                  </el-pagination>
                </div>
              </div>
            </div>
          </div>
          <div class="center_text" style="margin-top: 40px; height: 422px">
            <!--左上边框-->
            <div class="t_line_box">
              <i class="t_l_line"></i>
              <i class="l_t_line"></i>
            </div>
            <!--右上边框-->
            <div class="t_line_box">
              <i class="t_r_line"></i>
              <i class="r_t_line"></i>
            </div>
            <!--左下边框-->
            <div class="t_line_box">
              <i class="l_b_line"></i>
              <i class="b_l_line"></i>
            </div>
            <!--右下边框-->
            <div class="t_line_box">
              <i class="r_b_line"></i>
              <i class="b_r_line"></i>
            </div>
            <div class="main_title">Algorithm Management</div>
            <div class="bottom_center fl">
              <div class="main_table t_btn8">
                <!-- 查询 -->
                <div class="search">
                  <el-form
                    :inline="true"
                    :model="searchArea2"
                    class="demo-form-inline choose"
                  >
                    <el-form-item label="Algorithm Search:" style="margin-bottom: 5px">
                      <el-select
                        v-model="algorithmValue"
                        placeholder="Select Algorithm"
                        class="input"
                        @focus="getAlgorithmName()"
                      >
                        <el-option
                          v-for="item in algorithms"
                          :key="item.algorithmValue"
                          :label="item.algorithmValue"
                          :value="item.algorithmValue"
                        >
                        </el-option>
                      </el-select>
                    </el-form-item>
                    <el-form-item style="margin-left: 20px; margin-bottom: 5px">
                      <el-button type="primary" @click="searchAlgorithm()"
                        >Search</el-button
                      >
                    </el-form-item>
                    <el-form-item style="margin-left: 10px; margin-bottom: 5px">
                      <el-button type="warning" @click="examineBtn2()">Add</el-button>
                    </el-form-item>
                    <el-form-item style="margin-left: 10px; margin-bottom: 5px">
                      <el-button type="success" @click="refresh1()">Refresh</el-button>
                    </el-form-item>
                  </el-form>
                </div>
                <!-- 表格 -->
                <div class="box-content">
                  <el-table
                    :data="
                      tableData3.slice(
                        (currentPage4 - 1) * pageSize,
                        currentPage4 * pageSize
                      )
                    "
                    style="width: 100%"
                  >
                    <el-table-column
                      prop="name"
                      label="Algorithm Name"
                      style="width: 25%"
                    >
                    </el-table-column>
                    <el-table-column prop="mirror" label="Status" style="width: 25%">
                    </el-table-column>
                    <el-table-column prop="URL" label="URL" style="width: 25%">
                    </el-table-column>
                    <el-table-column prop="detail" label="Detail" style="width: 25%">
                    </el-table-column>
                  </el-table>
                </div>
                <!-- 页码 -->
                <div class="block">
                  <el-pagination
                    layout="total, prev, pager, next"
                    :total="pageSizes3"
                    :page-size="pageSize"
                    @current-change="getmirrorMessage()"
                    :current-page.sync="currentPage3"
                  >
                  </el-pagination>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <taskTransfer ref="msgBtn" v-if="Visiable"></taskTransfer>
    <taskPublish ref="msgBtn1" v-if="Visiable1"></taskPublish>
    <addAlgorithm ref="msgBtn2" v-if="Visiable2"></addAlgorithm>
    <warningMes ref="msgBtn4" v-if="Visiable4"></warningMes>
    <workLoad ref="msgBtn6" v-if="Visiable6"></workLoad>
  </div>
</template>
<script>
import taskTransfer from "../components/taskTransfer";
import addAlgorithm from "../components/addAlgorithm";
import taskPublish from "../components/taskPublish";
import warningMes from "../components/warningMes";
import workLoad from "../components/workflow.vue";
import * as echarts from "echarts";
import axios from "axios";
axios.defaults.baseURL = "/api";
export default {
  data() {
    return {
      // 与弹出窗口显示有关
      Visiable: false,
      Visiable1: false,
      Visiable2: false,
      Visiable4: false,
      Visiable6: false, // 工作流发布

      // 节点管理模块
      // 查询栏
      searchArea1: {},
      // 节点数据
      tableData1: [],
      // 获取节点名称
      nodes1: [],
      nodeValue1: "",
      // 关于节点的分页
      currentPage1: 1,
      pageSizes1: 0,
      pageSize: 4,

      // 任务管理模块
      // 查询栏
      searchArea2: {},
      // 当前节点
      currentNode: "",
      // 节点数据
      tableData2: [],
      // 获取节点名称
      nodes2: [],
      nodeValue2: "",
      // 关于节点的分页
      currentPage2: 1,
      pageSizes2: 0,

      // 资源管理模块
      nodecpuSources: [],
      nodememorySources: [],
      nodediskSources: [],

      // 关于节点的分页
      currentPage3: 1,
      pageSizes3: 0,

      // 关于算法的分页
      currentPage4: 1,
      pageSizes4: 0,
      // 算法数据
      tableData3: [],
      // 选择算法
      algorithms: [],
      algorithmValue: "",

      // 当前地址
      url: "",

      // 未读消息
      notReadMes: 0,
    };
  },
  mounted: function () {
    this.url = this.getFromSessionStorage("url");
    this.startCheckingConnection();
    this.getnodeMessage(),
      this.gettaskNum(),
      this.getcpuData(),
      this.getmemoryData(),
      this.getdiskData(),
      this.createSseConnect(2024),
      this.getAlgorithm();
  },
  beforeDestroy() {
    this.CloseSSE();
    this.stopCheckingConnection();
  },
  methods: {
    // 与弹窗有关
    examineBtn(data) {
      this.Visiable = true;
      this.$nextTick(() => {
        this.$refs.msgBtn.init(data);
      });
    },
    examineBtn1(data) {
      this.Visiable1 = true;
      this.$nextTick(() => {
        this.$refs.msgBtn1.init(data);
      });
    },
    examineBtn2(data) {
      this.Visiable2 = true;
      this.$nextTick(() => {
        this.$refs.msgBtn2.init(data);
      });
    },
    examineBtn4(data) {
      this.Visiable4 = true;
      this.$nextTick(() => {
        this.$refs.msgBtn4.init(data);
      });
      this.notReadMes = 0;
    },
    examineBtn6(data) {
      this.Visiable6 = true;
      this.$nextTick(() => {
        this.$refs.msgBtn6.init(data);
      });
    },

    // 获取会话内容
    getFromSessionStorage(key) {
      return window.sessionStorage.getItem(key);
    },

    // 节点管理模块
    // 搜索节点信息
    searchNode() {
      const that = this;
      let ipAddress = that.getFromSessionStorage("url");
      let nodeUrl = `http://${ipAddress}`;
      axios
        .post(
          `${nodeUrl}:5000/searchNode`,
          that.$qs.stringify({
            nodeName: that.nodeValue1,
          })
        )
        .then(function (response) {
          that.tableData1 = [];
          let data = response.data.nodeMessage;
          that.tableData1.push(data);
          that.pageSizes1 = 1;
          that.currentPage1 = 1;
        })
        .catch(function (error) {
          console.log(error);
        });
    },
    // 获取节点名称
    getnodeName1() {
      const that = this;
      let ipAddress = that.getFromSessionStorage("url");
      let nodeUrl = `http://${ipAddress}`;
      axios
        .get(`${nodeUrl}:5000/getnodeName`, {
          data: "",
        })
        .then(function (response) {
          let i = 0;
          that.nodes1 = [];
          for (i in response.data.node) {
            let node = response.data.node[i];
            that.nodes1.push({ nodeValue1: node, nodeLabel1: node });
          }
        })
        .catch(function (error) {
          console.log(error);
        });
    },
    // 请求节点信息
    getnodeMessage() {
      const that = this;
      let ipAddress = that.url;
      let nodeUrl = `http://${ipAddress}`;
      axios
        .post(
          `${nodeUrl}:5000/getnodeMessage`,
          that.$qs.stringify({
            currentPage: that.currentPage1,
            pageSize: that.pageSize,
          })
        )
        .then(function (response) {
          let i = 0;
          that.tableData1 = [];
          for (i in response.data.nodeMessage) {
            let data = response.data.nodeMessage[i];
            if (data["status"] == "Offline"){
              data["cpuSize"] = "-"
              data["memorySize"] = "-"
              data["diskSize"] = "-"
            }
            that.tableData1.push(data);
          }
          that.pageSizes1 = response.data.total;
        })
        .catch(function (error) {
          console.log(error);
        });
    },
    // 改变当前页数
    handleCurrentChange(val) {
      this.currentPage = val;
      this.getnodeMessage();
    },
    // 刷新
    refreshNode() {
      this.currentPage1 = 1;
      this.getnodeMessage();
    },

    // 任务管理模块
    // 搜索任务信息
    searchTask() {
      const that = this;
      let ipAddress = that.getFromSessionStorage("url");
      let nodeUrl = `http://${ipAddress}`;
      axios
        .post(
          `${nodeUrl}:5000/searchTask`,
          that.$qs.stringify({
            nodeName: that.nodeValue2,
            currentPage: that.currentPage2,
            pageSize: that.pageSize,
          })
        )
        .then(function (response) {
          that.currentNode = that.nodeValue2;
          that.tableData2 = [];
          let i = 0;
          for (i in response.data.taskMessage) {
            let data = response.data.taskMessage[i];
            that.tableData2.push(data);
          }
          that.pageSizes2 = response.data.total;
          console.log(that.tableData2);
        })
        .catch(function (error) {
          console.log(error);
        });
    },
    // 获取节点名称
    getnodeName2() {
      const that = this;
      let ipAddress = that.getFromSessionStorage("url");
      let nodeUrl = `http://${ipAddress}`;
      axios
        .get(`${nodeUrl}:5000/getnodeName`, {
          data: "",
        })
        .then(function (response) {
          let i = 0;
          that.nodes2 = [];
          for (i in response.data.node) {
            let node = response.data.node[i];
            that.nodes2.push({ nodeValue2: node, nodeLabel2: node });
          }
        })
        .catch(function (error) {
          console.log(error);
        });
    },
    // 删除任务确定窗口
    deleteWindow(data) {
      this.$confirm("Delete Task?", "Tip", {
        confirmButtonText: "Confirm",
        cancelButtonText: "Cancel",
        type: "warning",
      })
        .then(() => {
          this.deleteTask(data);
          this.$message({
            type: "success",
            message: "删除成功!",
          });
        })
        .catch(() => {});
    },
    // 删除任务
    deleteTask(data) {
      const that = this;
      let ipAddress = that.getFromSessionStorage("url");
      let nodeUrl = `http://${ipAddress}`;
      axios
        .post(
          `${nodeUrl}:5000/deleteTask`,
          that.$qs.stringify({
            taskID: data.taskID,
          })
        )
        .then(function (response) {
          console.log(response);
        })
        .catch(function (error) {
          console.log(error);
        });
    },

    // 任务分布
    // 获取任务数量
    gettaskNum() {
      const that = this;
      let nodeName = [];
      let nodeValue = [];
      const cpuBar = echarts.init(document.getElementById("chart_1"));
      let ipAddress = that.getFromSessionStorage("url");
      let nodeUrl = `http://${ipAddress}`;
      axios
        .post(
          `${nodeUrl}:5000/getbarData`,
          that.$qs.stringify({
            data: "taskNum",
          })
        )
        .then(function (response) {
          console.log(response);
          nodeName = response.data.nodeName;
          nodeValue = response.data.nodeValue;
          console.log(nodeName);
          let taskNums = [];
          let i = 0;
          for (i in nodeValue) {
            let taskNum = nodeValue[i];
            taskNums.push(taskNum);
          }
          cpuBar.setOption({
            tooltip: {
              trigger: "item",
              formatter: "{a} <br/>{b} : {c}",
            },
            xAxis: {
              type: "category",
              data: nodeName,
              axisLabel: {
                textStyle: {
                  color: "white", //修改坐标轴字体颜色
                },
                fontSize: 12, //调整坐标轴字体大小
                color: "white",
              },
            },
            yAxis: {
              name: "Task Num",
              nameTextStyle: {
                color: "white",
                fontSize: 13,
              },
              minInterval: 1,
              type: "value",
              axisLabel: {
                textStyle: {
                  color: "white", //修改坐标轴字体颜色
                },
                fontSize: 12, //调整坐标轴字体大小
                color: "white",
              },
            },
            series: [
              {
                name: "Task Num",
                type: "bar",
                data: taskNums,
                itemStyle: {
                  normal: {
                    //这里是重点
                    color: function (params) {
                      //注意，如果颜色太少的话，后面颜色不会自动循环，最好多定义几个颜色
                      var colorList = [
                        "#546fc6",
                        "#91cb74",
                        "#fac859",
                        "#ed6766",
                        "#72c0de",
                        "#749f83",
                        "#ca8622",
                      ];
                      return colorList[params.dataIndex % colorList.length];
                    },
                  },
                },
              },
            ],
          });
        })
        .catch(function (error) {
          console.log(error);
        });
    },

    // 资源占用
    // 颜色控制
    customColorMethod(percentage) {
      if (percentage < 21) {
        return "#f56c6c";
      } else if (percentage < 41) {
        return "#e6a23c";
      } else if (percentage < 61) {
        return "#5cb87a";
      } else if (percentage < 81) {
        return "#1989fa";
      } else {
        return "#6f7ad3";
      }
    },

    fractionToFloat(fractionString) {
      // 使用正则表达式匹配分子和分母
      const matches = fractionString.match(/^(\d+)\/(\d+)$/);
      if (!matches) {
        throw new Error("Invalid fraction string");
      }

      // 将匹配到的字符串转换为数字
      const numerator = parseInt(matches[1], 10);
      const denominator = parseInt(matches[2], 10);

      // 检查分母是否为零
      if (denominator === 0) {
        throw new Error("Division by zero");
      }

      // 执行除法并返回结果
      return numerator / denominator;
    },

    // cpu占用
    getcpuData() {
      const that = this;
      // 先在这里定义才能使用
      let nodeName = [];
      let nodeValue = [];
      let total = 0;
      let ipAddress = that.getFromSessionStorage("url");
      let nodeUrl = `http://${ipAddress}`;
      axios
        .post(
          `${nodeUrl}:5000/getbarData`,
          that.$qs.stringify({
            data: "cpuSize",
          })
        )
        .then(function (response) {
          console.log(response);
          nodeName = response.data.nodeName;
          nodeValue = response.data.nodeValue;
          total = response.data.total;
          that.nodecpuSources = [];
          let i = 0;
          for (i in nodeName) {
            let cpuData = {
              nodeName: nodeName[i],
              percentage: Math.round(that.fractionToFloat(nodeValue[i]) * 100),
            };
            that.nodecpuSources.push(cpuData);
          }
        })
        .catch(function (error) {
          console.log(error);
        });
    },
    // 获取内存占用
    getmemoryData() {
      const that = this;
      // 先在这里定义才能使用
      let nodeName = [];
      let nodeValue = [];
      let total = 0;
      let ipAddress = that.getFromSessionStorage("url");
      let nodeUrl = `http://${ipAddress}`;
      axios
        .post(
          `${nodeUrl}:5000/getbarData`,
          that.$qs.stringify({
            data: "memorySize",
          })
        )
        .then(function (response) {
          console.log(response);
          nodeName = response.data.nodeName;
          nodeValue = response.data.nodeValue;
          total = response.data.total;
          that.nodememorySources = [];
          let i = 0;
          for (i in nodeName) {
            let memoryData = {
              nodeName: nodeName[i],
              percentage: Math.round(that.fractionToFloat(nodeValue[i]) * 100),
            };
            that.nodememorySources.push(memoryData);
          }
        })
        .catch(function (error) {
          console.log(error);
        });
    },
    // 获取磁盘占用
    getdiskData() {
      const that = this;
      let nodeName = [];
      let nodeValue = [];
      let total = 0;
      let ipAddress = that.getFromSessionStorage("url");
      let nodeUrl = `http://${ipAddress}`;
      axios
        .post(
          `${nodeUrl}:5000/getbarData`,
          that.$qs.stringify({
            data: "diskSize",
          })
        )
        .then(function (response) {
          console.log(response);
          nodeName = response.data.nodeName;
          nodeValue = response.data.nodeValue;
          that.nodediskSources = [];
          let i = 0;
          for (i in nodeName) {
            let diskData = {
              nodeName: nodeName[i],
              percentage: Math.round(that.fractionToFloat(nodeValue[i]) * 100),
            };
            that.nodediskSources.push(diskData);
          }
        })
        .catch(function (error) {
          console.log(error);
        });
    },
    sourceRefresh() {
      this.getcpuData(), this.getmemoryData(), this.getdiskData();
    },

    // 启动检查连接的定时器
    startCheckingConnection() {
      this.timerId = setInterval(this.checkConnection, 5000); // 5000毫秒 = 5秒
      // 立即检查一次连接
      this.checkConnection();
    },
    // 停止检查连接的定时器
    stopCheckingConnection() {
      clearInterval(this.timerId);
    },
    // 检查连接状态(异步操作)
    async checkConnection() {
      try {
        let ipAddress = this.getFromSessionStorage("url");
        let nodeUrl = `http://${ipAddress}`;
        // 假设有一个API端点可以检查连接状态
        const response = await fetch(`${nodeUrl}:5000/check-connection`);
        console.log(response);
        if (!response.ok) {
          // 连接断开，弹出警告框并跳转到登录页面
          alert("Current Leader malfunction, please reconnect");
          this.$router.push("/"); // 假设登录页面的路由是'/login'
        }
      } catch (error) {
        console.log(error);
        // 网络请求失败，同样认为是连接断开
        alert("Current Leader malfunction, please reconnect！");
        this.$router.push("/"); // 跳转到登录页面
      }
    },
    // 建立SSE连接
    createSseConnect(clientid) {
      const that = this;
      if (window.EventSource) {
        let ipAddress = that.getFromSessionStorage("url");
        let nodeUrl = `http://${ipAddress}`;
        // 如果浏览器支持EventSource接口
        that.eventSource = new EventSource(
          `${nodeUrl}:5000/SSE_Connect` + "/" + clientid
        ); // 创建一个EventSource对象
        console.log(that.eventSource);

        that.eventSource.onmessage = (event) => {
          // 监听消息
          console.log("onmessage: " + event.data);
          let data = JSON.parse(event.data);
          console.log(data.date);
          that.$notify({
            title: "报警消息",
            message: "注意: " + data.type,
            type: "error",
          });
          that.notReadMes += 1;
        };

        that.eventSource.onopen = (event) => {
          // 建立连接时执行
          console.log("onopen: " + event);
        };

        that.eventSource.onerror = (event) => {
          // 处理错误
          console.log("onerror: " + event);
        };

        that.eventSource.close = (event) => {
          // 关闭时执行
          console.log("close: " + event);
        };
      } else {
        console.log("你的浏览器不支持SSE~");
      }
      console.log(" 测试 打印");
    },
    // 关闭SSE连接
    CloseSSE() {
      if (this.eventSource) {
        this.eventSource.close();
        this.eventSource = null;
      }
    },

    // 算法管理
    // 切换页码
    handleCurrentChange1(val) {
      this.currentPage4 = val;
      // this.getnodeMessage1()
      // this.generateRelate()
    },
    // 请求节点信息
    getAlgorithm() {
      const that = this;
      let ipAddress = that.getFromSessionStorage("url");
      let nodeUrl = `http://${ipAddress}`;
      axios
        .get(`${nodeUrl}:5000/getAlgorithm`, {
          data: "",
        })
        .then(function (response) {
          that.tableData3 = response.data.algorithms;
          that.pageSizes4 = response.data.pageSizes;
        })
        .catch(function (error) {
          console.log(error);
        });
    },
    refresh() {
      this.currentPage4 = 1;
      this.getAlgorithm();
    },
    getAlgorithmName() {
      const that = this;
      let ipAddress = that.getFromSessionStorage("url");
      let nodeUrl = `http://${ipAddress}`;
      axios
        .get(`${nodeUrl}:5000/getAlgorithmName`, {
          data: "",
        })
        .then(function (response) {
          let i = 0;
          that.algorithms = [];
          for (i in response.data.algorithms) {
            let algorithm = response.data.algorithms[i];
            that.algorithms.push({
              algorithmValue: algorithm,
              algorithmLabel: algorithm,
            });
          }
        })
        .catch(function (error) {
          console.log(error);
        });
    },
    searchAlgorithm() {
      const that = this;
      let ipAddress = that.getFromSessionStorage("url");
      let nodeUrl = `http://${ipAddress}`;
      axios
        .post(
          `${nodeUrl}:5000/searchAlgorithm`,
          that.$qs.stringify({
            algorithmName: that.algorithmValue,
          })
        )
        .then(function (response) {
          that.tableData3 = [];
          let data = response.data.algorithmMessage;
          that.tableData3.push(data);
          that.pageSizes4 = 1;
          that.currentPage4 = 1;
          that.$message({
            showClose: true,
            type: "success",
            message: "查询成功!",
          });
        })
        .catch(function (error) {
          console.log(error);
        });
    },
    refresh1() {
      this.currentPage4 = 1;
      this.getAlgorithm();
    },
  },
  components: {
    taskTransfer,
    taskPublish,
    addAlgorithm,
    warningMes,
    workLoad,
    // chooseNode
  },
};
</script>

<style scoped>
.header {
  position: relative;
}

.t_title {
  width: 100%;
  height: 100%;
  text-align: center;
  font-size: 2.5em;
  line-height: 80px;
  color: #fff;
}

#chart_map {
  cursor: pointer;
}

.search {
  font-size: 15px;
  font-weight: 600;
  color: #61d2f7;
}

.box-content {
  border: 1px;
  padding: 10px;
}

.choose {
  margin-top: 30px;
  font-size: 14px;
  font-weight: bolder;
  color: #555555;
  width: 100%;
  text-align: left;
  margin-left: 20px;
}

.input {
  width: 100%;
  margin-left: 20px;
}

.message {
  position: absolute;
  right: 5%;
  top: 22%;
  z-index: 2;
  color: #f56c6c;
  font-size: 30px;
}
</style>
