<template>
  <div>
    <H1>Profiles环境</H1>
    <div style="background:#eee;padding: 10px">
       <h3>
       <Tooltip content="应用列表" placement="top-start">
           <a href="#/applist">首页 </a>
        </Tooltip>
       > {{$route.params.app}}
       </h3>
    </div><br>
      <HandleModal ref="handleModal"></HandleModal>
    <br>
    <Table :columns="tableConf.columns" :data="tableData"></Table>

  </div>
</template>

<script>
  import Api from 'common/api';
  import HandleModal from './handleModal.vue';
  export default {
    name: 'ProfileList',
    data() {
      return {
        permission:-1,
        tableConf: {
          columns: [
            {
              title: '环境',
              key: 'profile',
              align: 'center',
              render: (h, param) => h('div', [
                h('a', {
                  domProps: {
                      href: '#/configs/'+ this.$route.params.app+'/'+param.row.profile,
                  },
                }, param.row.profile),
              ]),
            },
            {
              title:"描述",
              key:"name",
              align:"center"
            },
            {
              title: '操作',
              key: 'action',
              align: 'center',
              render: (h, param) => h('div', [
                h('Button', {
                  props: {
                    type: 'primary',
                    size: 'small',
                  },
                  style: {
                    marginRight: '10px',
                  },
                  on: {
                    click: () => {
                      this.modifyItem(param.row);
                    },
                  },
                }, '修改'),
                h('Button', {
                  props: {
                    type: 'error',
                    size: 'small',
                  },
                  style: {
                    marginRight: '10px',
                  },
                  attrs:{
                  	disabled: this.permission != 9,
                  },
                  on: {
                    click: () => {
                      this.delItem(param.row.app,param.row.profile);
                    },
                  },
                }, '删除'),
                h('Button', {
                  props: {
                    type: 'success',
                    size: 'small',
                
                  },
                  style: {
                    marginRight: '10px',
                  },
                  attrs:{
                  	disabled: !window.localStorage, 
                  },
                  on: {
                    click: () => {
                      this.copyProfile(param.row);
                    },
                  },
                }, '复制'),
              ]),
            },
          ],
          data: [],
        },
        code:0,
        msg:''
      };
    },
    components: {
      HandleModal,
    },
    computed: {
      tableData() {
        return this.$store.state.MOD.profileMsg.profileList;
      }
    },
    mounted() {
      const self = this;
      let app = self.$route.params.app;//从url中获取当前app
      self.$store.state.MOD.currentApp=app;
      self.$store.state.MOD.profileMsg={};
      self.$store.dispatch('getProfileList',function(){
        if(self.$store.state.MOD.profileMsg.code != 0){
          self.$Message.error("error：" + self.$store.state.MOD.profileMsg.msg);
        }
      });
      
         //查询当前用户的权限
      Api.queryCurUserAuth(this.$route.params.app).done((res) => {
          if(res && res.code == 0){
          	self.permission = res.data == 99?9 : res.data;
          }
      });
     
    },
    methods: {
      modifyItem(data) {
        this.$refs.handleModal.setModifyData(data);
      },
      copyProfile(data) {
        const json =  JSON.stringify(data) ;
        window.localStorage.setItem("cnt2-copy-profile",data.app +"@@"+data.profile);
        this.$Message.success('已复制到粘贴板');
        
      },
      delItem(app,profile) {

         this.$Modal.confirm({
                    title: '提示',
                    content: '<p>确认删除？</p>',
                    onOk: () => {
                        Api.delProfile(app,profile).done((res) => {
                          if(res.code==0){
                            this.$Message.success('删除成功');
                            this.$store.dispatch('getProfileList',app);
                          }else{
                             this.$Message.error("删除失败："+res.msg);
                          }
                        });
                    },
                    onCancel: () => {
                       
                    }
                });
      },
    },
  };
</script>

<style lang="scss">
</style>
