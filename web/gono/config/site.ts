export type SiteConfig = typeof siteConfig;

// siteConfig 网站配置项
export const siteConfig = {
  name: "Gono",
  description:
    "RBAC-with-Domains management website implemented by Next.js + HeroUI.",
  navMenuItems: [
    {
      label: "Profile",
      href: "/profile",
    },
    {
      label: "Settings",
      href: "/settings",
    },
    {
      label: "Logout",
      href: "/logout",
    },
  ],
  links: {
    github: "https://github.com/sfshf/gonoweb",
    docs: "https://heroui.com",
    signIn: "/sign-in",
  },
};

// siteMenus 网站菜单列表
export const siteMenus = [
  {
    label: "undefined",
    href: "/undefined",
  },
  {
    label: "user",
    href: "/user",
  },
  {
    label: "domain",
    href: "/domain",
  },
  {
    label: "role",
    href: "/role",
  },
  {
    label: "menu",
    href: "/menu",
  },
  {
    label: "api",
    href: "/api",
  },
];

// hasMenu 判断菜单是否存在
export const hasMenu = (
  identifier: string,
): { label: string; href: string } | null => {
  for (const item of siteMenus) {
    if (item.href == identifier) {
      return { ...item };
    }
  }
  return null;
};

// siteWidgets 网站控件列表
export const siteWidgets = {
  btn_add_user: "btn_add_user",
  btn_update_user: "btn_update_user",
  btn_delete_user: "btn_delete_user",
  btn_search_user: "btn_search_user",
  btn_add_domain: "btn_add_domain",
  btn_update_domain: "btn_update_domain",
  btn_delete_domain: "btn_delete_domain",
  btn_search_domain: "btn_search_domain",
  btn_add_role: "btn_add_role",
  btn_update_role: "btn_update_role",
  btn_delete_role: "btn_delete_role",
  btn_search_role: "btn_search_role",
  btn_add_menu: "btn_add_menu",
  btn_update_menu: "btn_update_menu",
  btn_delete_menu: "btn_delete_menu",
  btn_search_menu: "btn_search_menu",
  btn_add_widget: "btn_add_widget",
  btn_update_widget: "btn_update_widget",
  btn_delete_widget: "btn_delete_widget",
  btn_search_widget: "btn_search_widget",
  btn_add_api: "btn_add_api",
  btn_update_api: "btn_update_api",
  btn_delete_api: "btn_delete_api",
  btn_search_api: "btn_search_api",
};
