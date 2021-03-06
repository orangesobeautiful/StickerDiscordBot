import PageLayout from "layouts/PageLayout.vue";

const routes = [
  {
    path: "/",
    name: "home",
    component: PageLayout,
    children: [
      {
        path: "",
        component: () => import("pages/ImgManage.vue")
      },
      {
        path: "img-manage",
        component: () => import("pages/ImgManage.vue")
      },
      {
        path: "sticker-web-tutorial",
        component: () => import("pages/StickerWebTutorial.vue")
      },
      {
        path: "bot-instruction-tutorial",
        component: () => import("pages/Nothing.vue")
      }
    ]
  },
  {
    path: "/login",
    component: () => import("pages/Login.vue")
  },
  {
    path: "/success-logout",
    component: () => import("pages/SuccessLogout.vue")
  }
];

// Always leave this as last one
if (process.env.MODE !== "ssr") {
  routes.push({
    path: "*",
    component: () => import("pages/Error404.vue")
  });
}

export default routes;
