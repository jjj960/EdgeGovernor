<!-- 选择节点的窗口 -->
<template>
    <!--1.首先，弹窗页面中要有el-dialog组件即弹窗组件，我们把弹窗中的内容放在el-dialog组件中-->
    <!--2.设置:visible.sync属性，动态绑定一个布尔值，通过这个属性来控制弹窗是否弹出-->
    <el-dialog title="Select Node" :visible.sync="detailVisible" width="25%" style="top: 20%; position: absolute;">

        <el-row>
            <el-select v-model="nodeValue" placeholder="Select node" style="width: 100%;" @focus="getNode()">
                <el-option v-for="item in nodes" :key="item.nodeValue" :label="item.nodeLabel" :value="item.nodeValue">
                </el-option>
            </el-select>
        </el-row>

        <el-row style="text-align: right;">
            <el-button @click="disappear" style="margin-left: -80px; margin-top: 20px;">Cancel</el-button>
            <el-button type="primary" @click="onSubmit">Confirm</el-button>
        </el-row>

    </el-dialog>
</template>

<script>
import axios from 'axios'
axios.defaults.baseURL = '/api'
export default {
    name: "chooseNode",
    data() {
        return {
            detailVisible: false,
            // 选择节点
            nodes: [],
            nodeValue: ''
        }
    },
    methods: {
        //3.定义一个init函数，通过设置detailVisible值为true来让弹窗弹出，这个函数会在父组件的方法中被调用
        init(data) {
            this.detailVisible = true;
            //data是父组件弹窗传递过来的值，我们可以打印看看
            console.log(data);
        },
        disappear() {
            this.detailVisible = false;
        },
        // 从窗口向页面传值
        selectNodeOK() {
            this.$emit("select1", this.nodeValue)
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
        onSubmit() {
            if (this.nodeValue == '') {
                this.disappear()
                this.$message({
                    type: 'warning',
                    message: '字段不能为空'
                });
            } else {
                this.selectNodeOK()
                this.disappear()
                this.$message({
                    type: 'success',
                    message: '操作成功'
                });
            }
        }
    }
}
</script>
