<template>
    <div style="height: 100%;">
        <div class="wp">
            <div class="boardLogin">
                <div style="font-size: 50px; margin-bottom: 30px; color: white;">Welcome</div>
                <div>
                    <el-row style="margin-bottom: 20px; margin-top: 10px;">
                        <el-input v-model="loginName" placeholder="Enter username"></el-input>
                    </el-row>
                    <el-row style="margin-bottom: 20px;">
                        <el-input v-model="password" placeholder="Enter Password" show-password></el-input>
                    </el-row>
                    <el-row>
                        <el-button type="primary" @click="Login()" style="margin-top: 10px;">Login</el-button>
                    </el-row>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios'
export default {
    name: "Login",
    data() {
        return {
            loginName: "",
            password: ""
        }
    },
    methods: {
        saveToSessionStorage(key, value) {
            window.sessionStorage.setItem(key, value);
        },
        // 向后端发送用户名和密码
        Login() {
            const that = this
            let ipAddress = that.loginName
            let nodeUrl = `http://${ipAddress}`
            axios.post(`${nodeUrl}:5000/Login`, that.$qs.stringify({
                password: that.password
            })).then(function (response) {
                if (response.data.status == 'success') {
                    // 设置token令牌
                    that.saveToSessionStorage('token', "yes")
                    that.saveToSessionStorage('url', ipAddress)
                    that.$router.push({ name:'sourceSystem'})
                    that.$message({
                        message: '登录成功',
                        type: 'success'
                    });
                }
                else {
                    that.$message({
                        type: 'warning',
                        message: '登录失败'
                    });
                }
            })
                .catch(function (error) {
                    console.log(error)
                    that.$message({
                        type: 'warning',
                        message: '连接不上节点'
                    });
                })
        },
        goExample(){
            this.$router.push('/example')
            this.$message({
                type: 'success',
                message: '进入实例'
            })
        }
    }
}
</script>
