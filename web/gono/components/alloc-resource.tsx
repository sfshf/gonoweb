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
  Image,
  Selection,
} from "@heroui/react";
import { listRole } from "@/api/role";
import { useTranslation } from "react-i18next";
import { TDomain, TResource, TRole, TUser } from "@/zustand/types";
import { listResource } from "@/api/resource";
import { allocDomainRoleResources, domainRoleResources, domainRoles } from "@/api/casbin";

export interface AllocResourceModalProps {
  className?: string;
  classNames?: ModalProps["classNames"];
  domain: TDomain | null;
  isOpen: boolean;
  onClose: () => void;
  onOpenChange: () => void;
  isControlled?: boolean;
  getButtonProps?: (props?: any) => any;
  getDisclosureProps?: (props?: any) => any;
}

export const AllocResourceModal: FC<AllocResourceModalProps> = ({
  className,
  classNames,
  domain,
  isOpen,
  onClose,
  onOpenChange,
}) => {
  const { t } = useTranslation();
  const reducer = (state: any, action: any) => {
    switch (action.type) {
      case "roles":
        return { ...state, roles: action.value };
      case "menus":
        return { ...state, menus: action.value };
      case "widgets":
        return { ...state, widgets: action.value };
      case "apis":
        return { ...state, apis: action.value };
      case "ownedRoles":
        return { ...state, ownedRoles: action.value };
      case "selectedRole":
        return { ...state, selectedRole: action.value };
      case "ownedResources":
        return { ...state, ownedResources: action.value };
      case "selectedMenus":
        return { ...state, selectedMenus: action.value };
      case "selectedWidgets":
        return { ...state, selectedWidgets: action.value };
      case "selectedApis":
        return { ...state, selectedApis: action.value };
      default:
        return { ...state, ...action.value };
    }
  };
  const [state, dispatch] = React.useReducer(reducer, {
    roles: new Set(),
    menus: new Set(),
    widgets: new Set(),
    apis: new Set(),
    ownedRoles: new Set(),
    selectedRole: "",
    ownedResources: new Set(),
    selectedMenus: new Set(),
    selectedWidgets: new Set(),
    selectedApis: new Set(),
  });
  const effect = async () => {
    try {
      // 获取所有角色
      let roles: TRole[] = [];
      let resp: any = await listRole({
        page: 0,
        pageSize: 0,
        name: "",
        dxid: "",
      });
      if (resp && resp.data && resp.data.list) {
        roles = resp.data.list;
      }
      // 获取所有资源
      let menus: TResource[] = [];
      let widgets: TResource[] = [];
      let apis: TResource[] = [];
      resp = await listResource({
        page: 0,
        pageSize: 0,
        type: 0,
        name: "",
        identifier: "",
      });
      if (resp && resp.data && resp.data.list) {
        for (let i = 0; i < resp.data.list.length; i++) {
          switch (resp.data.list[i].type) {
            case 1: // menu
              menus.push(resp.data.list[i]);
              break;
            case 2: // widget
              widgets.push(resp.data.list[i]);
              break;
            case 3: // api
              apis.push(resp.data.list[i]);
              break;
          }
        }
      }
      // 获取该域租户下的已有角色的xid列表
      let ownedRoles: string[] = [];
      resp = await domainRoles(domain?.xid ?? "");
      if (resp && resp.data) {
        ownedRoles = resp.data;
      }
      dispatch({
        value: {
          roles: new Set(roles),
          menus: new Set(menus),
          widgets: new Set(widgets),
          apis: new Set(apis),
          ownedRoles: new Set(ownedRoles),
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

  const onPressRole = async (item: TRole) => {
    try {
      // 获取该域租户下的角色的资源identifier列表
      let ownedResources: string[] = [];
      let resp: any = await domainRoleResources(domain?.xid ?? "", item.xid);
      if (resp && resp.data) {
        ownedResources = resp.data;
      }
      let selectedMenus: string[] = [];
      let selectedWidgets: string[] = [];
      let selectedApis: string[] = [];
      for (let i = 0; i < ownedResources.length; i++) {
        for (const val of state.menus) {
          if (val.identifier === ownedResources[i]) {
            selectedMenus.push(ownedResources[i]);
          }
        }
        for (const val of state.widgets) {
          if (val.identifier === ownedResources[i]) {
            selectedWidgets.push(ownedResources[i]);
          }
        }
        for (const val of state.apis) {
          if (val.identifier === ownedResources[i]) {
            selectedApis.push(ownedResources[i]);
          }
        }
      }
      dispatch({
        value: {
          selectedRole: item.xid,
          ownedResources: new Set(ownedResources),
          selectedMenus: new Set(selectedMenus),
          selectedWidgets: new Set(selectedWidgets),
          selectedApis: new Set(selectedApis),
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
  const onPressConfirm = async () => {
    try {
      if (!domain || !state.selectedRole) {
        return;
      }
      const ownedResources = [
        ...state.selectedMenus,
        ...state.selectedWidgets,
        ...state.selectedApis,
      ];
      const resp: any = await allocDomainRoleResources({
        dxid: domain.xid,
        rxid: state.selectedRole,
        identifiers: ownedResources,
      });
      dispatch({
        value: {
          ownedResources: new Set(ownedResources),
        }
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
      size='full'
      isOpen={isOpen}
      placement='top-center'
      onOpenChange={onOpenChange}
      onClose={() => {
        dispatch({
          value: {
            domains: new Set(),
            roles: new Set(),
            menus: new Set(),
            widgets: new Set(),
            apis: new Set(),
            ownedRoles: new Set(),
            selectedRole: "",
            ownedResources: new Set(),
            selectedMenus: new Set(),
            selectedWidgets: new Set(),
            selectedApis: new Set(),
          },
        }); //清空状态
        onClose();
      }}
    >
      <ModalContent>
        {(onClose) => (
          <>
            <ModalHeader className='flex flex-col gap-1'>
              {t("alloc_resource.header")}
            </ModalHeader>
            <ModalBody className='flex flex-row'>
              {domain && (
                <>
                  <Listbox
                    classNames={{
                      base: "max-w-xs",
                      list: "max-h-[70vh] overflow-scroll",
                    }}
                    aria-label='roles'
                    variant='flat'
                    topContent={<div>{t("alloc_resource.label.roles")}</div>}
                    items={state.roles ?? []}
                    selectionMode='multiple'
                    selectedKeys={state.ownedRoles}
                  >
                    {(item: TRole) => (
                      <ListboxItem
                        key={item.xid}
                        textValue={item.name}
                        onPress={() => onPressRole(item)}
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
                  <Listbox
                    classNames={{
                      base: "max-w-xs",
                      list: "max-h-[70vh] overflow-scroll",
                    }}
                    aria-label='menus'
                    variant='flat'
                    topContent={<div>{t("alloc_resource.label.menus")}</div>}
                    items={state.menus ?? []}
                    selectionMode='multiple'
                    selectedKeys={state.selectedMenus}
                    onSelectionChange={(keys: Selection) => {
                      dispatch({
                        type: "selectedMenus",
                        value: keys === 'all' ? state.menus : keys,
                      });
                    }}
                  >
                    {(item: any) => (
                      <ListboxItem key={item.identifier} textValue={item.name}>
                        <div className='flex gap-2 items-center'>
                          <Image className='shrink-0' src={item.icon} />
                          <div className='flex flex-col'>
                            <span className='text-small'>{item.name}</span>
                            <span className='text-tiny text-default-400'>
                              {item.identifier}
                            </span>
                            <span className='text-tiny text-default-400'>
                              {item.intro}
                            </span>
                          </div>
                        </div>
                      </ListboxItem>
                    )}
                  </Listbox>
                  <Listbox
                    classNames={{
                      base: "max-w-xs",
                      list: "max-h-[70vh] overflow-scroll",
                    }}
                    selectionMode='multiple'
                    aria-label='widget'
                    variant='flat'
                    topContent={<div>{t("alloc_resource.label.widgets")}</div>}
                    items={state.widgets ?? []}
                    selectedKeys={state.selectedWidgets}
                    onSelectionChange={(keys: Selection) => {
                      dispatch({
                        type: "selectedWidgets",
                        value: keys === 'all' ? state.widgets : keys,
                      });
                    }}
                  >
                    {(item: any) => (
                      <ListboxItem key={item.identifier} textValue={item.name}>
                        <div className='flex gap-2 items-center'>
                          <Image className='shrink-0' src={item.icon} />
                          <div className='flex flex-col'>
                            <span className='text-small'>{item.name}</span>
                            <span className='text-tiny text-default-400'>
                              {item.identifier}
                            </span>
                            <span className='text-tiny text-default-400'>
                              {item.intro}
                            </span>
                          </div>
                        </div>
                      </ListboxItem>
                    )}
                  </Listbox>
                  <Listbox
                    classNames={{
                      base: "max-w-xs",
                      list: "max-h-[70vh] overflow-scroll",
                    }}
                    selectionMode='multiple'
                    aria-label='api'
                    variant='flat'
                    topContent={<div>{t("alloc_resource.label.apis")}</div>}
                    items={state.apis ?? []}
                    selectedKeys={state.selectedApis}
                    onSelectionChange={(keys: Selection) => {
                      dispatch({
                        type: "selectedApis",
                        value: keys === "all" ? state.apis : keys,
                      });
                    }}
                  >
                    {(item: any) => (
                      <ListboxItem key={item.identifier} textValue={item.name}>
                        <div className='flex gap-2 items-center'>
                          <Image className='shrink-0' src={item.icon} />
                          <div className='flex flex-col'>
                            <span className='text-small'>{item.name}</span>
                            <span className='text-tiny text-default-400'>
                              {item.identifier}
                            </span>
                            <span className='text-tiny text-default-400'>
                              {item.intro}
                            </span>
                          </div>
                        </div>
                      </ListboxItem>
                    )}
                  </Listbox>
                </>
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
                      menus: new Set(),
                      widgets: new Set(),
                      apis: new Set(),
                      ownedRoles: new Set(),
                      selectedRole: "",
                      ownedResources: new Set(),
                      selectedMenus: new Set(),
                      selectedWidgets: new Set(),
                      selectedApis: new Set(),
                    },
                  }); //清空状态
                  onClose();
                }}
              >
                {t("app.btn.close")}
              </Button>
              <Button
                disabled={!domain || !state.selectedRole}
                color={
                  !domain || !state.selectedRole
                    ? "default"
                    : "primary"
                }
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
