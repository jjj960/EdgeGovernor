<template>
    <div>
        <!-- style="position: absolute; top:0%" 很重要，可以解决下拉框偏移问题-->
        <el-dialog title="Task Deployment" :visible.sync="detailVisible" width="30%" style="position: absolute; top:0%">
            <el-row>
                <!-- 表单 -->
                <el-form :model="formInline" ref="form" label-width="174px">
                    <el-form-item label="Task Name:" class="choose itemCenter">
                        <el-input v-model="taskName" placeholder="Enter task name" class="input"></el-input>
                    </el-form-item>
                    <el-form-item label="Scheduler Algorithm:" class="choose itemCenter">
                        <el-select v-model="mirrorValue" placeholder="Select scheduler algorithm" class="input" @focus="getMirror()">
                            <el-option v-for="item in mirrors" :key="item.mirrorValue" :label="item.mirrorLabel"
                                :value="item.mirrorValue">
                            </el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="Image:" class="choose itemCenter">
                        <el-input v-model="mirrorName" placeholder="Enter task image" class="doubleInput"
                            style="margin-bottom: 10px;"></el-input>
                        <el-input v-model="mirrorVersion" placeholder="Enter image version" class="doubleInput"></el-input>
                    </el-form-item>
                    <el-form-item label="Deploy Mode:" class="choose itemCenter">
                        <el-select v-model="tpValue" placeholder="Select deployment time" class="doubleInput" style="margin-bottom: 10px;"
                            v-on:change="choosepostTime(tpValue)">
                            <el-option v-for="item in taskpostValue" :key="item.tpValue" :label="item.tpLabel"
                                :value="item.tpValue">
                            </el-option>
                        </el-select>
                        <el-input v-model="taskpostTime" class="doubleInput" :disabled="true"></el-input>
                    </el-form-item>
                    <el-form-item label="Request CPU(m)：" class="choose itemCenter">
                        <el-input v-model="cpuSize" placeholder="Enter Request CPU" class="input"></el-input>
                    </el-form-item>
                    <el-form-item label="Request Mem(Mb)：" class="choose itemCenter">
                        <el-input v-model="memorySize" placeholder="Enter Request Mem" class="input"></el-input>
                    </el-form-item>
                    <el-form-item label="Request Net:" class="choose itemCenter">
                        <el-input v-model="network" placeholder="Enter Request Net" class="input"></el-input>
                    </el-form-item>
                    <el-form-item label="Data Persistence：" class="choose itemCenter">
                        <el-select v-model="dataValue" placeholder="Please select" class="doubleInput"
                            style="margin-bottom: 10px;" v-on:change="chooseData(dataValue)">
                            <el-option v-for="item in dataSave" :key="item.dataValue" :label="item.dataLabel"
                                :value="item.dataValue">
                            </el-option>
                        </el-select>
                        <el-input v-model="diskSizeMB" class="doubleInput" :disabled="true"></el-input>
                    </el-form-item>
                    <el-form-item label="Deploy nodes:" class="choose itemCenter" style="text-align: center; width: 100%;">
                        <el-select v-model="nodeValue" placeholder="Please select" class="doubleInput"
                            style="margin-bottom: 10px;" v-on:change="chooseNode(nodeValue)">
                            <el-option v-for="item in node" :key="item.nodeValue" :label="item.nodeLabel"
                                :value="item.nodeValue">
                            </el-option>
                        </el-select>
                        <el-input v-model="nodeNumber" class="doubleInput" :disabled="true"></el-input>
                    </el-form-item>>
                </el-form>
            </el-row>
            <el-row style="text-align: right; background-color: #253a4f;">
                <el-button @click="disappear()" style="margin-left: -80px;">Cancel</el-button>
                <el-button type="primary" @click="taskPublish()">Publish</el-button>
            </el-row>
        </el-dialog>
        <!-- 窗口 -->
        <taskpost ref="msgBtn" v-if="Visiable" @select="selectTime"></taskpost>
        <chooseNode ref="msgBtn1" v-if="Visiable1" @select1="selectNode"></chooseNode>
    </div>
</template>

<script>
import taskpost from "./taskpost";
import filter from './filter.js';
import chooseNode from './chooseNode';
import axios from 'axios'
axios.defaults.baseURL = '/api'
export default {
    data() {
        return {
            detailVisible: false,

            // 与弹出窗口显示有关
            Visiable: false,
            Visiable1: false,

            taskName: '',
            mirrorName: '',
            mirrorVersion: '',
            cpuSize: '',
            // 内存大小
            memorySize: '',
            // 网络宽带
            network: '',
            // 任务发布模式下拉框
            taskpostValue: [{
                tpValue: 'scheRelease',
                tpLabel: 'Scheduled Deployment'
            }, {
                tpValue: 'insRelease',
                tpLabel: 'Single-task Deployment'
            }],
            tpValue: '',
            // 任务发布时间
            date1: '',
            date2: '',
            taskpostTime: '',
            // 是否数据持久化
            dataSave: [{
                dataValue: 'save',
                dataLabel: 'Yes'
            }, {
                dataValue: 'nosave',
                dataLabel: 'No'
            }],
            dataValue: '',
            // 请求磁盘大小
            diskSize: '',
            diskSizeMB: '',
            // 是否指定节点
            node: [{
                nodeValue: 'node',
                nodeLabel: 'Yes'
            }, {
                nodeValue: 'noNode',
                nodeLabel: 'No'
            }],
            nodeValue: '',
            // 节点号
            nodeNumber: '',
            // 获取镜像仓库
            mirrors: [],
            mirrorValue: ''
        }
    },
    methods: {
        init(data) {
            this.detailVisible = true;
            //data是父组件弹窗传递过来的值，我们可以打印看看
            console.log(data);
        },
        disappear() {
            this.detailVisible = false;
        },
        // 选择发布时间
        choosepostTime(tpValue) {
            if (tpValue == 'scheRelease') {
                this.examineBtn()
            } else if (tpValue == 'insRelease') {
                this.taskpostTime = ""
            }
        },
        examineBtn(data) {
            this.Visiable = true;
            this.$nextTick(() => {
                this.$refs.msgBtn.init(data);
            })
        },
        selectTime(date1, date2) {
            this.date1 = filter.formatDate(date1, 'yyyy-MM-dd')
            this.date2 = filter.formatDate(date2, 'HH:mm')
            this.taskpostTime = this.date1 + " " + this.date2
        },
        // 数据持久化
        chooseData(dataValue) {
            if (dataValue == 'save') {
                this.disksizeOpen()
            } else if (dataValue == 'nosave') {
                this.diskSize = ''
                this.diskSizeMB = ''
            }
        },
        disksizeOpen() {
            this.$prompt('', '请求磁盘大小(MB)', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                inputPlaceholder: '请输入请求磁盘大小'
            }).then(({ value }) => {
                if (value == '' || !(!isNaN(parseFloat(value)) && isFinite(value))) {
                    this.$message({
                        type: 'warning',
                        message: '字段输入有误'
                    });
                } else {
                    this.diskSize = value
                    this.diskSizeMB = value + "Mb"
                    this.$message({
                        type: 'success',
                        message: '操作成功'
                    });
                }
            });
        },
        // 节点选择
        chooseNode(nodeValue) {
            if (nodeValue == 'node') {
                this.examineBtn1()
            } else if (nodeValue == 'noNode') {
                this.nodeNumber = ""
            }
        },
        examineBtn1(data) {
            this.Visiable1 = true;
            this.$nextTick(() => {
                this.$refs.msgBtn1.init(data);
            })
        },
        selectNode(number) {
            this.nodeNumber = number
        },
        // 获取会话内容
        getFromSessionStorage(key) {
            return window.sessionStorage.getItem(key);
        },
        // 请求镜像仓库
        getMirror() {
            const that = this
            let ipAddress = that.getFromSessionStorage('url')
            let nodeUrl = `http://${ipAddress}`
            axios.get(`${nodeUrl}:5000/getAlgorithmName`, {
                data: ''
            }).then(function (response) {
                console.log(response)
                let i = 0
                that.mirrors = []
                for (i in response.data.algorithms) {
                    let mirror = response.data.algorithms[i]
                    that.mirrors.push({ mirrorValue: mirror, mirrorLabel: mirror })
                }
                console.log(that.mirrors)
            })
                .catch(function (error) {
                    console.log(error)
                })
        },
        // 提交任务
        taskPublish() {
            const that = this
            let ipAddress = that.getFromSessionStorage('url')
            let nodeUrl = `http://${ipAddress}`
            axios.post(`${nodeUrl}:5000/taskPublish`, that.$qs.stringify({
                taskName: that.taskName,
                mirrorValue: that.mirrorValue,
                mirrorName: that.mirrorName,
                mirrorVersion: that.mirrorVersion,
                // 任务发布模式，定时发的时间(taskpostTime1代表日期，taskpostTime2代表小时分钟)
                tpValue: that.tpValue,
                taskpostTime1: that.date1,
                taskpostTime2: that.date2,
                cpuSize: that.cpuSize,
                memorySize: that.memorySize,
                network: that.network,
                // 指示是否数据持久化，请求磁盘大小
                dataValue: that.dataValue,
                diskSize: that.diskSize,
                // 是否指定节点，节点编号
                nodeValue: that.nodeValue,
                nodeNumber: that.nodeNumber
            })).then(function (response) {
                that.$message({
                    type: 'success',
                    message: '操作成功'
                });
                console.log(response)
                that.disappear()
            })
                .catch(function (error) {
                    that.$message({
                        type: 'warning',
                        message: '操作失败'
                    });
                    console.log(error)
                })
        }
    },
    components: {
        taskpost,
        chooseNode
    }
}
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
</style>
