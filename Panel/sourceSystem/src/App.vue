<template>
    <div id="app">
        <router-view></router-view>
    </div>
</template>

<script>
// import leftNav from '@/components/common/leftNav'
export default {
    name: 'App',
    components: {
        // leftNav
    }
}
</script>

<style>
#app {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    color: #2c3e50;
    height: 100%;
}

.el-form-item__content {
    width: 61%;
}

.el-form-item__label {
  text-align: center;
}

.el-message-box__header {
    text-align: center;
}

.el-form-item__label {
    font-size: 15px;
    font-weight: 600;
    color: #61d2f7;
}

.el-input__inner {
    background-color: rgb(7, 41, 81);
    box-shadow: rgb(44, 88, 166) -10px 0px 15px inset, rgb(44, 88, 166) 10px 0px 15px inset;
    border-color: rgb(7, 41, 81);
    color: #C0C4CC;
}

.el-select-dropdown {
    background-color: rgb(7, 41, 81);
    box-shadow: rgb(44, 88, 166) -10px 0px 15px inset, rgb(44, 88, 166) 10px 0px 15px inset;
}

.el-table {
    background-color: transparent;
}

.el-table th.el-table__cell.is-leaf {
    background-color: transparent;
    color: #61d2f7;
}

.el-table tr {
    background-color: transparent;
    color: white;
    font-size: 14px;
}

.el-table td.el-table__cell,
.el-table th.el-table__cell.is-leaf {
    border-bottom: 0px;
}

.el-table--border::after,
.el-table--group::after,
.el-table::before {
    background: transparent;
}

.el-button--primary {
    color: #FFF;
    background-color: #218bb3;
    border-color: #218bb3;
}

.el-button--success {
    color: #FFF;
    background-color: #5fa350;
    border-color: #5fa350;
}

.el-button--warning {
    background-color: #ce9338;
    border-color: #ce9338;
    color: #FFF;
}

/* 鼠标经过颜色改为粉红色 */
.el-table__body tr:hover>td {
    background-color: rgb(7, 41, 81) !important;
}

.el-pager li{
    background: transparent;
    color: #555555;
}

.el-pager li.active {
    color: #61d2f7;
    cursor: default;
}

.el-pager li:hover {
    color: #61d2f7;
}

.el-pagination button:disabled{
    background: transparent;
    color: #555555;
}

.el-pagination .btn-next, .el-pagination .btn-prev {
    background: transparent;
    color: #555555;
}

.el-pagination button:hover{
    color: #61d2f7;
}

.el-pager li:hover{
    color: #61d2f7;
}

.el-dialog__header{
    background-color: #253a4f;
    border-radius: 21px 21px 0px 0px;
}

.el-dialog__body{
    background-color: #253a4f;
    border-radius: 0px 0px 21px 21px;
}

.el-dialog__title{
    color: white;
    font-weight: 600;
}

.el-input.is-disabled .el-input__inner{
    background-color: rgb(7, 41, 81);
    border: 0px;
}

.el-dialog__headerbtn{
    z-index: 9999;
    font-size: 19px;
}

.el-dialog{
    background-color: transparent;
}

.el-message-box{
    border-radius: 21px;
    background: rgb(7, 41, 81);
    border: 0px;
}

.el-message-box__title{
    color: white;
    font-weight: 600;
}

.el-message-box__message{
    color: #61d2f7;
}

.el-table__row{
    height: 47px;
}

.el-progress__text{
    vertical-align: top;
    color: white !important;
    width: 5%;
}

.el-progress-bar{
    margin-bottom: -10px !important;
}

/* 美化滚动条 */
.core_dialogue {
  height: 100%;
  overflow: scroll;
}
.core_dialogue::-webkit-scrollbar-track-piece {
  background-color: #f8f8f800;
}
.core_dialogue::-webkit-scrollbar {
  width: 6px;
  transition: all 2s;
}
.core_dialogue::-webkit-scrollbar-thumb {
  background-color: transparent;
  border-radius: 100px;
}
.core_dialogue::-webkit-scrollbar-thumb:hover {
  background-color: transparent;
}
.core_dialogue::-webkit-scrollbar-corner {
  background-color: rgba(255, 255, 255, 0);
}
/* 美化滚动条2 */
.el-table__body-wrapper::-webkit-scrollbar-track-piece {
  background-color: #f8f8f800;
}
.el-table__body-wrapper::-webkit-scrollbar {
  height: 8px;
  transition: all 2s;
}
.el-table__body-wrapper::-webkit-scrollbar-thumb {
  background-color: rgba(144,147,153,.3);
  border-radius: 100px;
}
.el-table__body-wrapper::-webkit-scrollbar-thumb:hover {
  background-color: rgba(144,147,153,.3);
}
.el-table__body-wrapper::-webkit-scrollbar-corner {
  background-color: rgba(255, 255, 255, 0);
}

.el-button--primary {
    color: #FFF;
    background-color: #218bb3;
    border-color: #218bb3;
}
.el-button--primary {
    color: #FFF;
    background-color: #409EFF;
    border-color: #409EFF;
}
.el-table--scrollable-x .el-table__body-wrapper {
    overflow-x: auto;
}
.el-upload__input {
    display: None!important;
}
.el-upload {
    display: inline-block;
}
.el-upload-dragger {
    background-color: transparent;
}
.el-upload-dragger .el-upload__text em {
    color: #61d2f7;
    font-style: normal;
}
.upload-demo {
    text-align: left;
}
.el-upload-list__item-name:hover{
    color:#409EFF;
}
.el-textarea__inner{
    background-color: rgb(7, 41, 81);
    -webkit-box-shadow: rgb(44, 88, 166) -10px 0px 15px inset, rgb(44, 88, 166) 10px 0px 15px inset;
    box-shadow: rgb(44, 88, 166) -10px 0px 15px inset, rgb(44, 88, 166) 10px 0px 15px inset;
    border-color: rgb(7, 41, 81);
    color: #C0C4CC;
    overflow: hidden;
}
.el-textarea.is-disabled .el-textarea__inner {
    background-color: rgb(7, 41, 81);
    border-color: rgb(7, 41, 81);
    cursor: not-allowed;
}
</style>
