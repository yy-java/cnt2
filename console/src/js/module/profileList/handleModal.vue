<template>
  <div>
    <Button type="success" @click="type=1;modal = true;formData={
          profile: '',
          name: ''
      }">添加profile
    </Button>
    <Modal
      :width="500"
      v-model="modal"
      :title="type==1? '新增Profile':'修改Profile'"
      ok-text="提交"
      :mask-closable="false"
      :footer-hide="true"
    >
      <Form :model="formData" ref="formData"   >
        <Form-item :rules="ruleKeyCustom" prop="profile" label="标记："required>
          <Input v-model="formData.profile"  :readonly="type==1?false:true" :disabled="type==1?false:true">
          </Input>
        </Form-item>
        <Form-item :rules="ruleCustom" prop="name" label="名称："required>
          <Input v-model="formData.name" >
          </Input>
        </Form-item>

        <Form-item >
          <Button type="primary" @click="submitData('formData',0)">
            {{type === 1 ? '保存' : '修改'}}
          </Button>
          &nbsp;&nbsp;&nbsp;&nbsp;
          <Button type="error" v-show="type == 1" :disabled="canPaste?false:true" :title="canPaste?'':'当前浏览器不支持'" @click="submitData('formData',1)">
              {{pasteBtnTitle}}
          </Button>
        </Form-item>
      </Form>

    </Modal>
  </div>
</template>

<script>
  import Api from 'common/api';
  import Check from 'common/check';
  export default {
    name: 'handleModal',
    data() {
      return {
        modal: false,
        loading: true,
        app:this.$route.params.app,
        type: 1, // 1:新增 2：修改
        canPaste: !!window.localStorage,
        pasteBtnTitle:'粘贴',
        formData: {
          profile: '',
          name: '',
          app:'',
          srcProfile:'',
        },
        ruleKeyCustom:  Check.ruleKeyCustom,
          ruleCustom:  Check.ruleCustom,
      };
    },
    updated:function(){
    	const self = this;
		self.pasteBtnTitle = function(){
			if(!window.localStorage)
				return "粘贴";
			var profiles = window.localStorage.getItem("cnt2-copy-profile");
			if(profiles){
		  		const profileArr = profiles.split("@@");
		  		if(self.app == profileArr[0]){
		  	    	return "保存并复制 App："+profileArr[0]+"，Profile："+profileArr[1]+" 的所有配置";
		  	    }
		  	}
		  	self.canPaste = false;
		  	return "粘贴板无信息";
	  	}();
	  if(this.type == 2){
	  	 this.$refs['']
	  }
      if(this.type==2){
         this.$refs['formData'].validate();
      }
    },
    methods: {
      setModifyData(data) {
        const formData = JSON.parse(JSON.stringify(data));
        this.formData = formData;
        this.modal = true;
        this.type = 2;
      },
      submitData(name,needCopy) {
        const self = this;
       	var srcProfile = '';
        if(needCopy){
        	var isError =false;
         	var profiles = window.localStorage.getItem("cnt2-copy-profile");
			if(!profiles){
	  	    	 isError = true
			}else{
		  		const profileArr = profiles.split("@@");
		  		if(self.app != profileArr[0] && needCopy){
		  	    	  isError = true
		  	    }else{
		  	    	srcProfile = profileArr[1];
		  	    }
	  	    }
	  	    if(isError){
	  	    	 this.$Message.error('无可复制的Profile');
	  	    	 return;
	  	    }
	  	}
        this.$refs[name].validate((valid) => {
            if (valid) {
                  self.formData.app = this.$route.params.app;
                  if (self.type === 1) {
                    self.formData.srcProfile = srcProfile;
                    Api.createProfile(self.formData).done((res) => {
                      if(res.code === 0){
                        this.$Message.success('添加成功');
                        this.modal=false;
                        self.$store.dispatch('getProfileList',self.formData.app);
                      }else{
                          this.$Message.error('添加失败，' + res.msg);
                      }
                    });
                  } else {
                    Api.updateProfile(self.formData).done((res) => {
                      if(res.code === 0){
                        this.$Message.success('修改成功');
                        this.modal=false;
                         self.$store.dispatch('getProfileList',self.formData.app);
                      }else{
                          this.$Message.error('修改失败，' + res.msg);
                      }
                    });
                  }
            }else {
                this.$Message.error('表单验证失败!');
            }
        });
      },
    },
  };
</script>

<style lang="scss">
</style>
