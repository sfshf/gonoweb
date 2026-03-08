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
  Select,
  SelectItem,
} from "@heroui/react";
import { listRole } from "@/api/role";
import { useTranslation } from "react-i18next";
import { TDomain, TResource, TRole } from "@/zustand/types";
import { listResource } from "@/api/resource";
import {
  allocDomainRoleResources,
  domainRoleResources,
  domainRoles,
} from "@/api/casbin";

type RoleT = TRole & {
  selected: boolean;
  owned: boolean;
};
type Auths = [] | ["read"] | ["write"];
type ResourceT = TResource & {
  auths: Auths;
  selected: boolean;
  owned: boolean;
};

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
  const auths = ["read", "write"];
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
      default:
        return { ...state, ...action.value };
    }
  };
  const [state, dispatch] = React.useReducer(reducer, {
    roles: [],
    menus: [],
    widgets: [],
    apis: [],
    ownedRoles: [],
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
      // 更新状态
      const rolets: RoleT[] = [];
      for (const r of roles) {
        let owned = false;
        for (const o of ownedRoles) {
          if (o === r.xid) {
            owned = true;
          }
        }
        rolets.push({
          ...r,
          selected: owned,
          owned,
        });
      }
      const menuts: ResourceT[] = [];
      for (const r of menus) {
        menuts.push({
          ...r,
          auths: [],
          selected: false,
          owned: false,
        });
      }
      const widgetts: ResourceT[] = [];
      for (const w of widgets) {
        widgetts.push({
          ...w,
          auths: [],
          selected: false,
          owned: false,
        });
      }
      const apits: ResourceT[] = [];
      for (const a of apis) {
        apits.push({
          ...a,
          auths: [],
          selected: false,
          owned: false,
        });
      }
      dispatch({
        value: {
          roles: rolets,
          menus: menuts,
          widgets: widgetts,
          apis: apits,
          ownedRoles: ownedRoles,
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
  const trimAct = (actObj: string): Auths => {
    if (actObj.startsWith("read ")) {
      return ["read"];
    } else if (actObj.startsWith("write ")) {
      return ["write"];
    }
    return [];
  };
  const onPressRole = async (item: RoleT) => {
    try {
      // 获取该域租户下的角色的资源identifier列表
      let ownedResources: string[] = [];
      let resp: any = await domainRoleResources(domain?.xid ?? "", item.xid);
      if (resp && resp.data) {
        ownedResources = resp.data ? resp.data : [];
      }
      const menus: ResourceT[] = JSON.parse(JSON.stringify(state.menus));
      menu: for (const val of menus) {
        // 菜单
        for (let i = 0; i < ownedResources.length; i++) {
          if (ownedResources[i].endsWith(val.identifier)) {
            val.owned = true;
            val.selected = true;
            val.auths = trimAct(ownedResources[i]);
            continue menu;
          }
        }
        val.owned = false;
        val.selected = false;
        val.auths = [];
      }
      const widgets: ResourceT[] = JSON.parse(JSON.stringify(state.widgets));
      widget: for (const val of widgets) {
        // 控件
        for (let i = 0; i < ownedResources.length; i++) {
          if (ownedResources[i].endsWith(val.identifier)) {
            val.owned = true;
            val.selected = true;
            val.auths = trimAct(ownedResources[i]);
            continue widget;
          }
        }
        val.owned = false;
        val.selected = false;
        val.auths = [];
      }
      const apis: ResourceT[] = JSON.parse(JSON.stringify(state.apis));
      api: for (const val of apis) {
        // API
        for (let i = 0; i < ownedResources.length; i++) {
          if (val.identifier === ownedResources[i]) {
            val.owned = true;
            val.selected = true;
            val.auths = [];
            continue api;
          }
        }
        val.owned = false;
        val.selected = false;
        val.auths = [];
      }
      const roles: RoleT[] = JSON.parse(JSON.stringify(state.roles));
      for (const val of roles) {
        if (val.xid === item.xid) {
          val.selected = true; // 当前role被选中
        } else {
          val.selected = false; // 其他role的selected状态被清空
        }
      }
      // 更新状态
      dispatch({
        value: {
          roles: roles,
          menus: menus,
          widgets: widgets,
          apis: apis,
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
      if (!domain) {
        throw "null domain";
      }
      // selected role
      let rxid = "";
      for (const r of state.roles) {
        if (r.selected) {
          if (rxid) {
            throw "multiple roles are selected";
          }
          rxid = r.xid;
        }
      }
      const identifiers: string[] = [];
      for (const m of state.menus) {
        if (m.selected && m.auths[0]) {
          identifiers.push(m.auths[0] + " " + m.identifier);
        }
      }
      for (const m of state.widgets) {
        if (m.selected && m.auths[0]) {
          identifiers.push(m.auths[0] + " " + m.identifier);
        }
      }
      for (const m of state.apis) {
        if (m.selected) {
          identifiers.push(m.identifier);
        }
      }
      const resp: any = await allocDomainRoleResources({
        dxid: domain.xid,
        rxid,
        identifiers,
      });
      const roles: RoleT[] = JSON.parse(JSON.stringify(state.roles));
      for (const val of roles) {
        if (val.xid === rxid) {
          val.owned = true; // 当前role被拥有
        }
      }
      // 更新状态
      dispatch({
        value: {
          roles: roles,
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
      size='full'
      isOpen={isOpen}
      placement='top-center'
      onOpenChange={onOpenChange}
      onClose={() => {
        dispatch({
          value: {
            roles: [],
            menus: [],
            widgets: [],
            apis: [],
            ownedRoles: [],
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
              <Listbox
                classNames={{
                  base: "max-w-xs",
                  list: "max-h-[70vh] overflow-scroll",
                }}
                aria-label='roles'
                variant='flat'
                topContent={<div>{t("alloc_resource.label.roles")}</div>}
                items={state.roles}
              >
                {(item: RoleT) => (
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
                items={state.menus}
              >
                {(item: ResourceT) => (
                  <ListboxItem
                    key={item.identifier}
                    textValue={item.name}
                    endContent={
                      <Select
                        size='sm'
                        className='h-[20px] w-[85px]'
                        placeholder={t("alloc_resource.placeholder.none")}
                        onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
                          const menus = JSON.parse(JSON.stringify(state.menus));
                          for (const m of menus) {
                            if (m.identifier === item.identifier) {
                              m.auths = e.target.value
                                ? ([e.target.value] as Auths)
                                : undefined;
                              m.selected = e.target.value ? true : false;
                              break;
                            }
                          }
                          dispatch({ type: "menus", value: menus });
                        }}
                        selectedKeys={item.auths}
                      >
                        {auths.map((rw) => (
                          <SelectItem key={rw}>
                            {t("alloc_resource.label." + rw)}
                          </SelectItem>
                        ))}
                      </Select>
                    }
                  >
                    <div className='flex gap-2 items-center'>
                      {/* <Image className='shrink-0' src={item.icon} /> */}
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
                aria-label='widget'
                variant='flat'
                topContent={<div>{t("alloc_resource.label.widgets")}</div>}
                items={state.widgets}
              >
                {(item: ResourceT) => (
                  <ListboxItem
                    key={item.identifier}
                    textValue={item.name}
                    endContent={
                      <Select
                        size='sm'
                        className='h-[20px] w-[85px]'
                        placeholder={t("alloc_resource.placeholder.none")}
                        onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
                          const widgets = JSON.parse(
                            JSON.stringify(state.widgets),
                          );
                          for (const m of widgets) {
                            if (m.identifier === item.identifier) {
                              m.auths = e.target.value
                                ? ([e.target.value] as Auths)
                                : undefined;
                              m.selected = e.target.value ? true : false;
                              break;
                            }
                          }
                          dispatch({ type: "widgets", value: widgets });
                        }}
                        selectedKeys={item.auths}
                      >
                        {auths.map((rw) => (
                          <SelectItem key={rw}>
                            {t("alloc_resource.label." + rw)}
                          </SelectItem>
                        ))}
                      </Select>
                    }
                  >
                    <div className='flex gap-2 items-center'>
                      {/* <Image className='shrink-0' src={item.icon} /> */}
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
                aria-label='api'
                variant='flat'
                topContent={<div>{t("alloc_resource.label.apis")}</div>}
                items={state.apis}
                selectionMode='multiple'
                selectedKeys={state.apis.map((item: ResourceT) =>
                  item.selected ? item.identifier : "",
                )}
                onSelectionChange={(keys: Selection) => {
                  const selects = Array.from(new Set(keys));
                  const apis: ResourceT[] = JSON.parse(
                    JSON.stringify(state.apis),
                  );
                  if (keys === "all") {
                    for (const a of apis) {
                      a.selected = true;
                    }
                  } else {
                    api: for (const a of apis) {
                      for (const key of selects) {
                        if (a.identifier === key) {
                          a.selected = true;
                          continue api;
                        }
                      }
                      a.selected = false;
                    }
                  }
                  dispatch({
                    type: "apis",
                    value: apis,
                  });
                }}
              >
                {(item: ResourceT) => (
                  <ListboxItem key={item.identifier} textValue={item.name}>
                    <div className='flex gap-2 items-center'>
                      {/* <Image className='shrink-0' src={item.icon} /> */}
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
            </ModalBody>
            <ModalFooter>
              <Button
                color='danger'
                variant='flat'
                onPress={() => {
                  dispatch({
                    value: {
                      roles: [],
                      menus: [],
                      widgets: [],
                      apis: [],
                      ownedRoles: [],
                    },
                  }); //清空状态
                  onClose();
                }}
              >
                {t("app.btn.close")}
              </Button>
              <Button color='primary' onPress={onPressConfirm}>
                {t("app.btn.confirm")}
              </Button>
            </ModalFooter>
          </>
        )}
      </ModalContent>
    </Modal>
  );
};
