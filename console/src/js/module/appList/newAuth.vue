<template>
  <div>
    <Button type="primary" @click="type=1;newAuthModal = true;newAuthformData={
          uid: '',
          username: '',
          permission: '1',
          app:''
      }">
   新增 </Button>
    <Modal
      :width="300"
      v-model="newAuthModal"
      :transfer="false"
      :title="type==1?'用户授权':'修改授权'"
      ok-text="提交"
      :mask-closable="false"
      :footer-hide="true">
      <Form :model="newAuthformData" ref="newAuthformData">
        <Form-item :rules="ruleNumberCustom"  prop="uid" label="uid：" required>
          <Input v-model="newAuthformData.uid"  :readonly="type==1?false:true" :disabled="type==1?false:true">
          </Input>
        </Form-item>
        <Form-item :rules="ruleCustom" prop="username" label="用户名：" required>
          <Input v-model="newAuthformData.username"   :readonly="type==1?false:true">
          </Input>
        </Form-item>
         <Form-item  prop="permission" label="授权类型：" required>
          <RadioGroup v-model="newAuthformData.permission">
              <Radio label="1">
                  <span>开发</span>
              </Radio>
              <Radio label="9">
                  <span>管理</span>
              </Radio>
          </RadioGroup>
        </Form-item>
        
        <Form-item>
          <Button type="primary" @click="submitData('newAuthformData')">
            {{type === 1 ? '提交' : '修改'}}
          </Button>
        </Form-item>
      </Form>

    </Modal>
  </div>
</template>

<script>
  import Api from 'common/api';
  import UserAuth from './userAuth.vue';
   import Check from 'common/check';
  export default {
    name: 'newAuth',
    data() {
     
      return {
        newAuthModal: false,
        loading: true,
        type: 1, // 1:新增 2：修改
        newAuthformData: {
          uid: '',
          username: '',
          permission: '',
          app:''
        },
        ruleCustom:  Check.ruleCustom,
        ruleNumberCustom: Check.ruleNumberCustom, 
      };
    },
     updated:function(){
      if(this.type==2){
         this.$refs['newAuthformData'].validate();
      }
    },
    methods: {
      setModifyData(data) {
        const formData = JSON.parse(JSON.stringify(data));  
        this.newAuthformData = formData;
        this.newAuthModal = true;
        this.type = 2;
      },
      submitData(name) {
        this.$refs[name].validate((valid) => {
            if (valid) {
                const self = this;
                if (self.type === 1) {
                   self.newAuthformData.app=this.$store.state.MOD.currentApp;      
                   Api.createAuth(self.newAuthformData).done((res) => {
                    
                    if(res.code === 0){      
                      this.$Message.success('添加成功');
                      this.newAuthModal=false;
                      // location.reload();
                      this.$store.dispatch('getUserAuthList');
                       this.$store.dispatch('getTableList');
                    }else{
                       this.$Message.error('添加失败，' + res.msg);
                    }
                  });
                } else {
                  Api.updateAuth(self.newAuthformData).done((res) => {
                   
                    if(res.code === 0){
                      this.$store.dispatch('getUserAuthList');
                      this.$Message.success('修改成功');
                      this.newAuthModal=false;
                      this.$store.dispatch('getUserAuthList');
                       this.$store.dispatch('getTableList');
                    }else{
                        this.$Message.error('修改失败，' + res.msg);
                    }
                  });
                }
            } else {
                this.$Message.error('表单验证失败!');
            }
        });

      },
    },
  };
</script>

<style lang="scss">
</style>