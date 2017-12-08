import * as types from '../mutation-types';
import Api from '../../common/api';

// initial state
const states = {
  currentUid:'',
  tableList: [],
  profileMsg:{},
  configKVList:[],//全部配置内容
  waitApproveConfig:[],//待审核
  waitPublishConfig:[],//待发布
  configHistoryList:[],
  currentApp:'',
  currentProfile:'',
  currentUserList:[]
};

// getters
const getters = {};

// actions
const actions = {
  getTableList({
                 commit,
               }) {
    Api.getAppList().done((res) => {
      
      commit(types.GET_LIST, {
        list: res.data?res.data:[],
      });
    });
  },

  getProfileList({
                 commit,
               },callback){
      states.profileMsg={};
       Api.profiles(states.currentApp).done(
       (res) => {
	       commit(types.GET_PROFILE_LIST, {
	        data: res.data,
	        code:res.code,
	        msg:res.msg,
          callback:callback
	      })
    });    
  },
    getConfigList({
                 commit,
               }) {
    Api.queryConfig(states.currentApp,states.currentProfile).done((res) => {
      commit(types.GET_CONFIG_KV_LIST, {
        list: res.data,
        
      });
    });
  },
  getConfigHistoryList({
                 commit,
               },config) {
    Api.queryConfigHistory(config).done((res) => {
      commit(types.GET_CONFIG_HISTORY_LIST, {
        list: res.data,
      });
    });
  },
  getUserAuthList({
                 commit,
               }) {
    Api.getUserListByApp(states.currentApp).done((res) => {
      commit(types.GET_CURRENT_USER_LIST, {
        list: res.data,
      });
    });
  },
}
// mutations
const mutations = {
  [types.GET_LIST](state, data) {
    const nowState = state;
    nowState.tableList = data.list;
  },
  

  [types.GET_PROFILE_LIST](state, info) {
    const nowState = state;
    nowState.profileMsg={};//注意清空，每个app的profile不同
    nowState.profileMsg.profileList=info.data;
    nowState.profileMsg.code =info.code;
    nowState.profileMsg.msg=info.msg;
    
    
  },
  // 获取配置list后过滤分配到待审核和待发布
  [types.GET_CONFIG_KV_LIST](state, data) {
    const nowState = state;
    // nowState.configKVList = data.list;
    let configKVList=[];
    let waitApproveConfig = [];
    let waitPublishConfig = [];
    for(var d in data.list){
    	  configKVList.push(data.list[d]);
          if(data.list[d].version <= data.list[d].publishedVersion){
              continue;
          }
          if(data.list[d].approveType==1){
            waitApproveConfig.push(data.list[d]);
          }else if(data.list[d].approveType==2){
            waitPublishConfig.push(data.list[d]);
          }
    }
    nowState.configKVList=configKVList;
    nowState.waitApproveConfig=waitApproveConfig;
    nowState.waitPublishConfig=waitPublishConfig;
  },
  [types.GET_CONFIG_HISTORY_LIST](state, data) {
    const nowState = state;
    nowState.configHistoryList = data.list;
  },
 [types.GET_CURRENT_USER_LIST](state, data) {
    const nowState = state;
    let tempList=[];
    for(let i=0;i<data.list.length;i++){
      //过滤掉自己
      if(data.list[i].uid.toString()==nowState.currentUid){
          continue;
      }
     data.list[i].permissionStr=(data.list[i].permission==9)?'管理':(data.list[i].permission==1)?'开发':'';
     tempList.push(data.list[i]);
    }
    nowState.currentUserList = tempList;
  },
};

export default {
  state: states,
  getters,
  actions,
  mutations,
};
