<!-- 选择节点的窗口 -->
<template>
    <!--1.首先，弹窗页面中要有el-dialog组件即弹窗组件，我们把弹窗中的内容放在el-dialog组件中-->
    <!--2.设置:visible.sync属性，动态绑定一个布尔值，通过这个属性来控制弹窗是否弹出-->
    <div>
        <el-dialog title="Task Migration" :visible.sync="detailVisible" width="25%" style="top: 10%; position: absolute;">
            <el-row style="background-color: #253a4f;">
                <el-form ref="form" :model="form" label-width="80px">
                    <el-form-item label="ID:">
                        <el-input v-model="form.taskID" disabled></el-input>
                    </el-form-item>
                    <el-form-item label="Task Name:">
                        <el-input v-model="form.taskName" disabled></el-input>
                    </el-form-item>
                    <el-form-item label="From:">
                        <el-input v-model="form.nodeName" disabled></el-input>
                    </el-form-item>
                    <el-form-item label="To:">
                        <el-select v-model="nodeValue" placeholder="Select Node" style="width: 100%;" @focus="getNode()">
                            <el-option v-for="item in nodes" :key="item.nodeValue" :label="item.nodeLabel"
                                :value="item.nodeValue">
                            </el-option>
                        </el-select>
                    </el-form-item>
                </el-form>
            </el-row>

            <el-row style="text-align: right; background-color: #253a4f;">
                <el-button @click="disappear" style="margin-left: -80px; margin-top: 20px;">Cancel</el-button>
                <el-button type="primary" @click="postTransfer()">OK</el-button>
            </el-row>

        </el-dialog>
    </div>
</template>

<script>
import axios from 'axios'
import filter from './filter.js';
axios.defaults.baseURL = '/api'
export default {
    name: "taskTransfer",
    data() {
        return {
            detailVisible: false,
            // 表单
            form: {
                nodeName: "",
                taskID: "",
                taskName: ""
            },
            // 选择节点
            nodes: [],
            nodeValue: '',
            // CPU、内存、磁盘大小
            cpuSize: '',
            memorySize: '',
            diskSize: ''
        }
    },
    methods: {
        //3.定义一个init函数，通过设置detailVisible值为true来让弹窗弹出，这个函数会在父组件的方法中被调用
        init(data) {
            this.detailVisible = true;
            this.form.nodeName = data.nodeName
            this.form.taskID = data.taskID
            this.form.taskName = data.taskName
            this.cpuSize = data.cpuSize
            this.memorySize = data.memorySize
            this.diskSize = data.diskSize
            //data是父组件弹窗传递过来的值，我们可以打印看看
            console.log(data);
        },
        disappear() {
            this.detailVisible = false;
        },
        // 获取会话内容
        getFromSessionStorage(key) {
            return window.sessionStorage.getItem(key);
        },
        // 向后端发送请求获取节点
        getNode() {
            const that = this
            let ipAddress = that.getFromSessionStorage('url')
            let nodeUrl = `http://${ipAddress}`
            axios.get(`${nodeUrl}:5000/getnodeNameAb`, {
                data: ''
            }).then(function (response) {
                console.log(response)
                let i = 0
                that.nodes = []
                for (i in response.data.node) {
                    let node = response.data.node[i]
                    that.nodes.push({ nodeValue: node, nodeLabel: node })
                }
                console.log(that.nodes)
            })
                .catch(function (error) {
                    console.log(error)
                })
        },
        // 向后端发送转移命令
        postTransfer() {
            const that = this
            let date = new Date()
            let date1 = filter.formatDate(date, 'yyyy-MM-dd')
            let date2 = filter.formatDate(date, 'HH:mm')
            let ipAddress = that.getFromSessionStorage('url')
            let nodeUrl = `http://${ipAddress}`
            axios.post(`${nodeUrl}:5000/postTransfer`, that.$qs.stringify({
                taskID: that.form.taskID,
                taskName: that.form.taskName,
                From: that.form.nodeName,
                To: that.nodeValue,
                cpuSize: that.cpuSize,
                memorySize: that.memorySize,
                diskSize: that.diskSize,
                date1: date1,
                date2: date2
            })).then(function (response) {
                console.log(response)
                that.disappear()
                that.$message({
                    type: 'success',
                    message: '操作成功'
                });
            })
                .catch(function (error) {
                    console.log(error)
                    that.disappear()
                    that.$message({
                        type: 'warning',
                        message: '操作失败'
                    });
                })
        }
    }
}
</script>

