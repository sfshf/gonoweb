import { create } from "zustand";
import { persist } from "zustand/middleware";
import { AuthStore, TDomain, TResource, TRole, TUser } from "./types";

export const LocalStorageKey_AuthStore = "user_auth";
export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get, store) => ({
      token: null,
      user: null,
      domain: null,
      role: null,
      menus: null,
      widgets: null,
      setAuth: (
        token: string | null,
        user: TUser | null,
        domain: TDomain | null,
        role: TRole | null,
        menus: TResource[] | null,
        widgets: TResource[] | null,
      ) => set({ token, user, domain, role, menus, widgets }),
      clearAuth: () => {
        useAuthStore.persist.clearStorage();
        set(store.getInitialState());
      },
    }),
    {
      name: LocalStorageKey_AuthStore, // storage key
      // 默认使用 localStorage
    },
  ),
);

export const isRoot = (): boolean => {
  const user = useAuthStore.getState().user;
  if (user && user.email === process.env.NEXT_PUBLIC_ROOT_EMAIL) {
    return true;
  }
  return false;
};

// hasWidget 判断用户有没有控件权限
export const hasWidget = (identifier: string): boolean => {
  const widgets = useAuthStore.getState().widgets;
  if (!widgets) {
    return false;
  }
  for (const item of widgets) {
    if (item.identifier == identifier) {
      return true;
    }
  }
  return false;
};
