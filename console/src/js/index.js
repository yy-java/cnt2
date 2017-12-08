import VueRouter from 'vendor/vue-router';
import iView from 'iview';
import 'iview/dist/styles/iview.css';
import routers from './router/index';
import store from './store/index';
import App from './module/app.vue';
import User from './lib/user';

var init = function () {
  const router = new VueRouter({
    routes: routers,
  });
  Vue.use(iView);
  new Vue({ // eslint-disable-line
    el: '#app',
    router,
    store,
    components: {
      App,
    },
    template: '<App/>',
  });
};

function isIE() {
  var rVal;
  if (window.navigator.userAgent.indexOf('MSIE') >= 1) {
    rVal = true;
  } else {
    rVal = false;
  }
  return rVal;
}

if (isIE()) {
  	alert('不支持当前浏览器，请使用谷歌浏览器、360浏览器、QQ浏览器'); // eslint-disable-line
  	location.href = 'http://se.360.cn/';
} else if (!User.isLogin()) {
  	User.showLoginBox(location.href.replace(/#.*$/,''));
} else {
  	init();
}
 Date.prototype.Format = function(fmt)
	{ //author: meizz
	  var o = {
	    "M+" : this.getMonth()+1,                 //月份
	    "d+" : this.getDate(),                    //日
	    "h+" : this.getHours(),                   //小时
	    "m+" : this.getMinutes(),                 //分
	    "s+" : this.getSeconds(),                 //秒
	    "q+" : Math.floor((this.getMonth()+3)/3), //季度
	    "S"  : this.getMilliseconds()             //毫秒
	  };
	  if(/(y+)/.test(fmt))
	    fmt=fmt.replace(RegExp.$1, (this.getFullYear()+"").substr(4 - RegExp.$1.length));
	  for(var k in o)
	    if(new RegExp("("+ k +")").test(fmt))
	  fmt = fmt.replace(RegExp.$1, (RegExp.$1.length==1) ? (o[k]) : (("00"+ o[k]).substr((""+ o[k]).length)));
	  return fmt;
	}
