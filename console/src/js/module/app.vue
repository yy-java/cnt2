
<template>
  <div class="layout">
    <div class="layout-ceiling">
      <div class="layout-logo">
        配置中心
      </div>
      <div class="layout-ceiling-main" style='color:white'>
         {{getUsername()}}  <a  @click="signOut">【退出】</a>
      </div>
    </div>
    <div class="content">
      <router-view></router-view>
    </div>
  </div>
</template>
<script>
var Util = require('../lib/util');
import APi from 'common/api';

export default {
    name: 'app',
    methods: {
      getUsername(){
      	  this.$Message.config({top: 100})
          this.$store.state.MOD.currentUid=Util.getCookie("yyuid");
          return Util.getCookie("username");
      },
      signOut(){
          APi.logOut().done((res)=>{

            Util.clearCookie();
             location.reload();
           });
      },
    }
};

</script>


<style lang="scss">
  .content {
    padding: 10px;
  }

  .layout {
    border: 1px solid #d7dde4;
    background: #f5f7f9;
    position: relative;
    border-radius: 4px;
    overflow: hidden;
  }

  .layout-logo {
    width: 100px;
    height: 30px;
    line-height: 30px;
    text-align: center;
    /*background: #5b6270;*/
    border-radius: 3px;
    float: left;
    position: relative;
    left: 20px;
    color: #fff;
  }

  .layout-header {
    height: 60px;
    background: #fff;
    box-shadow: 0 1px 1px rgba(0, 0, 0, .1);
  }

  .layout-copy {
    text-align: center;
    padding: 10px 0 20px;
    color: #9ea7b4;
  }

  .layout-ceiling {
    background: #464c5b;
    padding: 10px 0;
    overflow: hidden;
  }

  .layout-ceiling-main {
    float: right;
    margin-right: 15px;
  }

  .layout-ceiling-main a {
    color: #9ba7b5;
  }

</style>
