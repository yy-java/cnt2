import AppList from 'module/appList/index.vue';
import ProfileList from 'module/profileList/index.vue';
import Configs from 'module/config/index.vue';

const routes = [
  {
    path: '/applist',
    component: AppList,
  },
  {
    path: '/profiles/:app',
    component: ProfileList,
  },
  {
    path: '/configs/:app/:profile',
    component: Configs,
  },
  {
    path: '*',
    redirect: '/applist',
  },
];

export default routes;
