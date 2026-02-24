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
import { editRole, listRole } from "@/api/role";
import { useTranslation } from "react-i18next";
import { TDomain, TResource, TRole, TUser } from "@/zustand/types";
import { listResource } from "@/api/resource";
import { listDomain } from "@/api/domain";
import { allocRoleInDomain, domainRoleResources, domainRoles } from "@/api/casbin";

export interface AllocRoleModalProps {
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

export const AllocRoleModal: FC<AllocRoleModalProps> = ({
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
      case "selectedDomain":
        return { ...state, selectedDomain: action.value };
      case "ownedRoles":
        return { ...state, ownedRoles: action.value };
      case "selectedRoles":
        return { ...state, selectedRoles: action.value };
      default:
        return { ...state, ...action.value };
    }
  };
  const [state, dispatch] = React.useReducer(reducer, {
    domains: new Set(),
    roles: new Set(),
    selectedDomain: "",
    ownedRoles: new Set(),
    selectedRoles: new Set(),
  });
  const effect = async () => {
    try {
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
      dispatch({
        value: {
          domains: new Set(domains),
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
  const onPressDomain = async (item: TDomain) => {
    try {
      // 获取该域租户下的角色的xid列表
      let ownedRoles: string[] = [];
      let resp: any = await domainRoles(item.xid);
      if (resp && resp.data) {
        ownedRoles = resp.data;
      }
      dispatch({
        value: {
          selectedDomain: item.xid,
          selectedRoles: new Set(ownedRoles),
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
  const onPressConfirm = async () => {
    try {
      if (!user || !state.selectedDomain) {
        return;
      }
      const resp: any = await allocRoleInDomain({
        xid: user.xid,
        dxid: state.selectedDomain,
        rxids: Array.from(state.selectedRoles),
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
              {state.selectedDomain && (
                <Listbox
                  classNames={{
                    base: "max-w-xs",
                    list: "max-h-[70vh] overflow-scroll",
                  }}
                  aria-label='roles'
                  variant='flat'
                  topContent={<div>{t("alloc_role.label.roles")}</div>}
                  items={state.roles ?? []}
                  selectionMode='multiple'
                  selectedKeys={state.selectedRoles}
                  onSelectionChange={(keys: Selection) => {
                    dispatch({
                      type: "selectedRoles",
                      value: keys === 'all' ? state.roles : keys,
                    });
                  }}
                >
                  {(item: TRole) => (
                    <ListboxItem
                      key={item.xid}
                      textValue={item.name}
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
                      selectedDomain: "",
                      ownedRoles: new Set(),
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
                color={
                  !user || !state.selectedDomain
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
