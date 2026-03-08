"use client";

import { FC } from "react";
import { ModalProps } from "@heroui/modal";
import * as React from "react";
import {
  addToast,
  Button,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Listbox,
  ListboxItem,
  Selection,
} from "@heroui/react";
import { listRole } from "@/api/role";
import { useTranslation } from "react-i18next";
import { AuthStore, TDomain, TResource, TRole, TUser } from "@/zustand/types";
import { listDomain } from "@/api/domain";
import { userDomains, userRolesInDomain } from "@/api/casbin";
import { switchRole } from "@/api/user";
import { useAuthStore } from "@/zustand/user";

export interface SwitchRolesModalProps {
  className?: string;
  classNames?: ModalProps["classNames"];
  user: TUser | null;
  isOpen: boolean;
  onClose: () => void;
  onOpenChange: () => void;
  isControlled?: boolean;
  getButtonProps?: (props?: any) => any;
  getDisclosureProps?: (props?: any) => any;
}

export const SwitchRolesModal: FC<SwitchRolesModalProps> = ({
  className,
  classNames,
  user,
  isOpen,
  onClose,
  onOpenChange,
}) => {
  const { t } = useTranslation();
  const reducer = (state: any, action: any) => {
    switch (action.type) {
      case "domains":
        return { ...state, domains: action.value };
      case "roles":
        return { ...state, roles: action.value };
      case "ownedRoles":
        return { ...state, ownedRoles: action.value };
      case "selectedDomain":
        return { ...state, selectedDomain: action.value };
      case "selectedRoles":
        return { ...state, selectedRoles: action.value };
      default:
        return { ...state, ...action.value };
    }
  };
  const [state, dispatch] = React.useReducer(reducer, {
    domains: new Set(),
    roles: new Set(),
    ownedRoles: new Set(),
    selectedDomain: "",
    selectedRoles: new Set(),
  });
  const effect = async () => {
    try {
      if (!user) {
        throw t("app.prompt.noSignIn");
      }
      // 获取所有域租户
      let domains: TDomain[] = [];
      let resp: any = await listDomain({
        page: 0,
        pageSize: 0,
        name: "",
      });
      if (resp && resp.data && resp.data.list) {
        domains = resp.data.list;
      }
      // 获取所有角色
      let roles: TRole[] = [];
      resp = await listRole({
        page: 0,
        pageSize: 0,
        name: "",
        dxid: "",
      });
      if (resp && resp.data && resp.data.list) {
        roles = resp.data.list;
      }
      // 获取用户被分配到的域租户列表
      let ownedDxids: string[] = [];
      resp = await userDomains(user.xid);
      if (resp && resp.data) {
        ownedDxids = resp.data;
      }
      // 过滤域租户列表
      let ownedDomains: TDomain[] = [];
      for (let i = 0; i < domains.length; i++) {
        if (ownedDxids.includes(domains[i].xid)) {
          ownedDomains.push(domains[i]);
        }
      }
      dispatch({
        value: {
          domains: new Set(ownedDomains),
          roles: new Set(roles),
        },
      });
    } catch (e) {
      addToast({
        title: (e as Error).cause as string,
        description: (e as Error).message,
        color: "danger",
      });
    }
  };
  const role = useAuthStore((state: AuthStore) => state.role);
  const onPressDomain = async (item: TDomain) => {
    try {
      if (!user) {
        throw t("app.prompt.noSignIn");
      }
      // 获取该域租户下的角色的xid列表
      let ownedRxids: string[] = [];
      let resp: any = await userRolesInDomain(user.xid, item.xid);
      if (resp && resp.data) {
        ownedRxids = resp.data;
      }
      // 过滤角色列表
      let ownedRoles: TRole[] = [];
      for (let role of state.roles) {
        if (ownedRxids.includes(role.xid)) {
          ownedRoles.push(role);
        }
      }
      dispatch({
        value: {
          selectedDomain: item.xid,
          ownedRoles: new Set(ownedRoles),
          selectedRoles: new Set([role?.xid]),
        },
      });
    } catch (e) {
      addToast({
        title: (e as Error).cause as string,
        description: (e as Error).message,
        color: "danger",
      });
    }
  };
  // 存储用户信息 -- 当前的角色和域租户，能访问到的菜单和控件
  const setAuth = useAuthStore((state: AuthStore) => state.setAuth);
  const onPressConfirm = async () => {
    try {
      if (!user) {
        throw t("app.prompt.noSignIn");
      }
      if (!state.selectedDomain) {
        return t("switch_roles.prompt.noSelectedDomain");
      }
      if (!state.selectedRoles || state.selectedRoles.size == 0) {
        return t("switch_roles.prompt.noSelectedRole");
      }
      const resp: any = await switchRole(
        state.selectedDomain,
        state.selectedRoles.values().next().value,
      );
      setAuth(
        resp.data.token,
        resp.data.user,
        resp.data.domain,
        resp.data.role,
        resp.data.menus,
        resp.data.widgets,
      );
      // 关闭Modal
      onClose();
      // 还原
      dispatch({
        value: {
          domains: new Set(),
          roles: new Set(),
          ownedRoles: new Set(),
          selectedDomain: "",
          selectedRoles: new Set(),
        },
      });
      addToast({
        title: t("app.prompt.ok"),
        description: resp.msg,
        color: "success",
      });
    } catch (e) {
      addToast({
        title: (e as Error).cause as string,
        description: (e as Error).message,
        color: "danger",
      });
    }
  };
  React.useEffect(() => {
    if (isOpen) {
      effect();
    }
  }, [isOpen]);
  return (
    <Modal
      size='xl'
      isOpen={isOpen}
      placement='top-center'
      onOpenChange={onOpenChange}
      onClose={() => {
        dispatch({
          value: {
            domains: new Set(),
            roles: new Set(),
            selectedDomain: "",
            ownedRoles: new Set(),
            selectedRoles: new Set(),
          },
        }); //清空状态
        onClose();
      }}
    >
      <ModalContent>
        {(onClose) => (
          <>
            <ModalHeader className='flex flex-col gap-1'>
              {t("alloc_role.header")}
            </ModalHeader>
            <ModalBody className='flex flex-row'>
              <Listbox
                classNames={{
                  base: "max-w-xs",
                  list: "max-h-[70vh] overflow-scroll",
                }}
                aria-label='domains'
                variant='flat'
                topContent={<div>{t("alloc_role.label.domains")}</div>}
                items={state.domains ?? []}
              >
                {(item: TDomain) => (
                  <ListboxItem
                    key={item.xid}
                    textValue={item.name}
                    onPress={() => onPressDomain(item)}
                  >
                    <div className='flex gap-2 items-center'>
                      <div className='flex flex-col'>
                        <span className='text-small'>{item.name}</span>
                        <span className='text-tiny text-default-400'>
                          {item.intro}
                        </span>
                      </div>
                    </div>
                  </ListboxItem>
                )}
              </Listbox>
              {state.selectedDomain && state.ownedRoles && (
                <Listbox
                  classNames={{
                    base: "max-w-xs",
                    list: "max-h-[70vh] overflow-scroll",
                  }}
                  aria-label='roles'
                  variant='flat'
                  topContent={<div>{t("alloc_role.label.roles")}</div>}
                  items={state.ownedRoles}
                  selectionMode='single'
                  selectedKeys={state.selectedRoles}
                  onSelectionChange={(keys: Selection) => {
                    dispatch({
                      type: "selectedRoles",
                      value:
                        keys === "all"
                          ? state.ownedRoles.map((r: TRole) => r.xid)
                          : keys,
                    });
                  }}
                >
                  {(item: TRole) => (
                    <ListboxItem key={item.xid} textValue={item.name}>
                      <div className='flex gap-2 items-center'>
                        <div className='flex flex-col'>
                          <span className='text-small'>{item.name}</span>
                          <span className='text-tiny text-default-400'>
                            {item.intro}
                          </span>
                        </div>
                      </div>
                    </ListboxItem>
                  )}
                </Listbox>
              )}
            </ModalBody>
            <ModalFooter>
              <Button
                color='danger'
                variant='flat'
                onPress={() => {
                  dispatch({
                    value: {
                      domains: new Set(),
                      roles: new Set(),
                      ownedRoles: new Set(),
                      selectedDomain: "",
                      selectedRoles: new Set(),
                    },
                  }); //清空状态
                  onClose();
                }}
              >
                {t("app.btn.close")}
              </Button>
              <Button
                disabled={!user || !state.selectedDomain}
                color={!user || !state.selectedDomain ? "default" : "primary"}
                onPress={onPressConfirm}
              >
                {t("app.btn.confirm")}
              </Button>
            </ModalFooter>
          </>
        )}
      </ModalContent>
    </Modal>
  );
};
