<template>
  <div> 
    <Modal
      :width="1000"
      v-model="modal" 
      :mask-closable="false"
      :footer-hide="true"
    >  
    <p slot="header" style="color:#f60;text-align:center">
          
            <span>{{currentApp}} 其他用户列表</span>
    </p>
    <NewAuth ref="newAuth" ></NewAuth><br>
    <Table :columns="tableUser.columns" :data="tableData"></Table>
    </Modal>
  </div>
</template>
<script>
  import Api from 'common/api';
  import NewAuth from './newAuth.vue';
  export default {
    name: 'authModal',
    data() {
      return {
        currentApp:'',
        tableUser:{
          columns:[ 
           {
              title: 'Uid',
              key: 'uid',
              align: 'center',
            },
            {
              title: '用户名',
              key: 'uname',
              align: 'center',
            },
            {
              title: '权限',
              key: 'permissionStr',
              align: 'center',
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
                    marginRight: '5px',
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
                  on: {
                    click: () => {
                      this.delItem(param.row.id);
                    },
                  },
                }, '删除'),
              ]),
            },
            ],
            data:[]
        },
        modal: false,
        loading: true,
        formData: {
          app: '',
          name: '',
          appType: '',
          chargerUid: '',
          charger: '',
        } 
      };
    },
    components: {
      
      NewAuth
    },
    computed: {
       tableData() {
        return this.$store.state.MOD.currentUserList;
      },
    },
    mounted() {
      const self = this;
     
    },
    methods: {
      setModifyData(data) {
        data.permission=permissionStr=="管理"?9:1;
        const formData = JSON.parse(JSON.stringify(data));
        this.formData = formData;
        this.modal = true;
        this.type = 2;
      },
      getUserList(){  
        this.modal = true;
        this.currentApp=this.$store.state.MOD.currentApp;
        this.$store.state.MOD.currentUserList=[];
        this.$store.dispatch('getUserAuthList');    
      },
      modifyItem(data) {    
        // data.permission=data.permission=='管理'?9:1;   
        data.username=data.uname
        this.$refs.newAuth.setModifyData(data);
      },
      delItem(id){  

        this.$Modal.confirm({
                    title: '提示',
                    content: '<p>确认删除？</p>',
                    onOk: () => {
                        Api.delAuth(id).done((res) =>{
                          if(res.code==0){
                            // alert("删除成功");
                            this.$Message.success('删除成功');
                            this.getUserList();
                          }else{
                             this.$Message.error("删除失败："+res.msg);
                          }
                        });
                    },
                    onCancel: () => {
                       
                    }
                });
       
      
      }
    },
  };
</script>

<style lang="scss">
</style>
