export type SiteConfig = typeof siteConfig;

export const siteConfig = {
  name: "Gono",
  description:
    "RBAC-with-Domains management website implemented by Next.js + HeroUI.",
  navItems: [
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
  ],
  navMenuItems: [
    {
      label: "Profile",
      href: "/profile",
    },
    {
      label: "Dashboard",
      href: "/dashboard",
    },
    {
      label: "Projects",
      href: "/projects",
    },
    {
      label: "Team",
      href: "/team",
    },
    {
      label: "Calendar",
      href: "/calendar",
    },
    {
      label: "Settings",
      href: "/settings",
    },
    {
      label: "Help & Feedback",
      href: "/help-feedback",
    },
    {
      label: "Logout",
      href: "/logout",
    },
  ],
  links: {
    github: "https://github.com/heroui-inc/heroui",
    twitter: "https://twitter.com/hero_ui",
    docs: "https://heroui.com",
    discord: "https://discord.gg/9b6yyZKmH4",
    signIn: "/sign-in",
  },
};

export const hasMenu = (
  identifier: string
): { label: string; href: string } | null => {
  for (const item of siteConfig.navItems) {
    if (item.href == identifier) {
      return { ...item };
    }
  }
  return null;
};
