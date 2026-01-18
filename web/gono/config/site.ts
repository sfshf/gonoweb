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
    label: "Undefined",
    href: "/undefined",
  },
  {
    label: "User",
    href: "/user",
  },
  {
    label: "Domain",
    href: "/domain",
  },
  {
    label: "Role",
    href: "/role",
  },
  {
    label: "Menu",
    href: "/menu",
  },
  {
    label: "API",
    href: "/api",
  },
];

// hasMenu 判断菜单是否存在
export const hasMenu = (
  identifier: string
): { label: string; href: string } | null => {
  for (const item of siteMenus) {
    if (item.href == identifier) {
      return { ...item };
    }
  }
  return null;
};

// siteWidgets 网站控件列表
export const siteWidgets = [
  {
    id: "btn_add_user",
    label: "",
  },
  {
    id: "btn_update_user",
    label: "",
  },
  {
    id: "btn_delete_user",
    label: "",
  },
  {
    id: "btn_search_user",
    label: "",
  },
  {
    id: "btn_add_domain",
    label: "",
  },
  {
    id: "btn_update_domain",
    label: "",
  },
  {
    id: "btn_delete_domain",
    label: "",
  },
  {
    id: "btn_search_domain",
    label: "",
  },
  {
    id: "btn_add_role",
    label: "",
  },
  {
    id: "btn_update_role",
    label: "",
  },
  {
    id: "btn_delete_role",
    label: "",
  },
  {
    id: "btn_search_role",
    label: "",
  },
  {
    id: "btn_add_menu",
    label: "",
  },
  {
    id: "btn_update_menu",
    label: "",
  },
  {
    id: "btn_delete_menu",
    label: "",
  },
  {
    id: "btn_search_menu",
    label: "",
  },
  {
    id: "btn_add_widget",
    label: "",
  },
  {
    id: "btn_update_widget",
    label: "",
  },
  {
    id: "btn_delete_widget",
    label: "",
  },
  {
    id: "btn_search_widget",
    label: "",
  },
  {
    id: "btn_add_api",
    label: "",
  },
  {
    id: "btn_update_api",
    label: "",
  },
  {
    id: "btn_delete_api",
    label: "",
  },
  {
    id: "btn_search_api",
    label: "",
  },
];

// hasWidgets 判断控件是否存在
export const hasWidgets = (
  identifier: string
): { id: string; label: string } | null => {
  for (const item of siteWidgets) {
    if (item.id == identifier) {
      return { ...item };
    }
  }
  return null;
};
