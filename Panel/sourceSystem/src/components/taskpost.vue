<!-- 选择任务发布时间的窗口 -->
<template>
    <!--1.首先，弹窗页面中要有el-dialog组件即弹窗组件，我们把弹窗中的内容放在el-dialog组件中-->
    <!--2.设置:visible.sync属性，动态绑定一个布尔值，通过这个属性来控制弹窗是否弹出-->
    <el-dialog title="Set Deployment Time" :visible.sync="detailVisible" width="25%" style="top: 20%; position: absolute;">
        <el-row>
            <el-col :span="11" style="margin-left: 2%;">
                <el-date-picker type="date" placeholder="Select Date" v-model="postTime.date1"
                    style="width: 100%;" :clearable="false"></el-date-picker>
            </el-col>
            <el-col :span="11" style="margin-left: 4%;">
                <el-time-picker placeholder="Select Time" v-model="postTime.date2" style="width: 100%;" :clearable="false"></el-time-picker>
            </el-col>
        </el-row>
        <el-row style="text-align: right;">
            <el-button @click="disappear" style="margin-left: -80px; margin-top: 20px;">Cancel</el-button>
            <el-button type="primary" @click="onSubmit">OK</el-button>
        </el-row>
    </el-dialog>
</template>

<script>
export default {
    name: "taskpost",
    data() {
        return {
            detailVisible: false,
            postTime: {
                name: '',
                region: '',
                date1: '',
                date2: '',
                delivery: false,
                type: [],
                resource: '',
                desc: ''
            },
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
        selectTimeOK() {
            this.$emit("select", this.postTime.date1, this.postTime.date2)
        },
        onSubmit() {
            if (this.postTime.date1 == '' || this.postTime.date2 == '') {
                this.disappear()
                this.$message({
                    type: 'warning',
                    message: '字段不能为空'
                });
            } else {
                this.selectTimeOK()
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
<style scoped>
.el-form-item__content {
    width: 100% !important;
}
</style>
