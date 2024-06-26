<!-- 选择节点的窗口 -->
<template>
    <!--1.首先，弹窗页面中要有el-dialog组件即弹窗组件，我们把弹窗中的内容放在el-dialog组件中-->
    <!--2.设置:visible.sync属性，动态绑定一个布尔值，通过这个属性来控制弹窗是否弹出-->
    <div>
        <el-dialog title="Add Algorithm" :visible.sync="detailVisible" width="30%" style="top: 5%; position: absolute;">
            <el-row style="background-color: #253a4f;">
                <el-form ref="form" :model="form" label-width="160px">
                    <el-form-item label="Algorithm Name:">
                        <el-input v-model="form.name" placeholder="Please enter a name for the algorithm"></el-input>
                    </el-form-item>
                    <el-form-item label="URL:">
                        <el-input v-model="form.url" placeholder="Please enter the URL address of the algorithm"></el-input>
                    </el-form-item>
                    <el-form-item label="Detail:">
                        <el-input v-model="form.detail" placeholder="Please enter the details of the algorithm"></el-input>
                    </el-form-item>
                    <el-form-item label="Status:">
                        <el-select v-model="useValue" placeholder="Whether the algorithm is enabled" style="width: 100%;">
                            <el-option v-for="item in uses" :key="item.useValue" :label="item.useLabel"
                                :value="item.useValue">
                            </el-option>
                        </el-select>
                    </el-form-item>
                </el-form>
            </el-row>

            <el-row style="text-align: right; background-color: #253a4f;">
                <el-button @click="disappear" style="margin-left: -80px; margin-top: 20px;">取消</el-button>
                <el-button type="primary" @click="addMirror()">确定</el-button>
            </el-row>

        </el-dialog>
    </div>
</template>
    
<script>
import axios from 'axios'
axios.defaults.baseURL = '/api'
export default {
    name: "addAlgorithm",
    data() {
        return {
            detailVisible: false,
            // 表单
            form: {
                name: "",
                url: "",
                detail: ""
            },
            // 下拉框
            uses: [{
                useValue: 'Enable',
                useLabel: 'Enable'
            }, {
                useValue: 'Deable',
                useLabel: 'Deable'
            }],
            useValue: '',
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
        // 向后端发送添加命令
        addMirror() {
            const that = this
            axios.post('/addAlgorithm', that.$qs.stringify({
                name: that.form.name,
                url: that.form.url,
                detail: that.form.detail,
                use: that.useValue,
            })).then(function (response) {
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

    