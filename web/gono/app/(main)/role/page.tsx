"use client";
import * as React from "react";
import {
  Table,
  TableHeader,
  TableColumn,
  TableBody,
  TableRow,
  TableCell,
  getKeyValue,
  addToast,
  Button,
  Input,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  useDisclosure,
} from "@heroui/react";
import { listRole, addRole, editRole, deleteRole } from "@/api/role";
import { Pagination } from "@heroui/pagination";
import { useTranslation } from "react-i18next";
import {
  FunnelIcon,
  PencilSquareIcon,
  PlusIcon,
  DocumentTextIcon,
  TrashIcon,
  UserIcon,
} from "@/components/icons";
import { hasWidget, useAuthStore } from "@/zustand/user";
import { AuthStore, TRole } from "@/zustand/types";
import { siteWidgets } from "@/config/site";
import { useSearchParams } from "next/navigation";

const AddRole = ({
  isOpen,
  onClose,
  onOpenChange,
}: {
  isOpen: boolean;
  onClose: () => void;
  onOpenChange: () => void;
  isControlled?: boolean;
  getButtonProps?: (props?: any) => any;
  getDisclosureProps?: (props?: any) => any;
}) => {
  const { t } = useTranslation();
  const reducer = (state: any, action: any) => {
    switch (action.type) {
      case "name":
        return { ...state, name: action.value };
      case "intro":
        return { ...state, intro: action.value };
      default:
        return { ...state, ...action.value };
    }
  };
  const [state, dispatch] = React.useReducer(reducer, {
    name: "",
    intro: "",
  });
  React.useEffect(() => {
    dispatch({ value: { name: "", intro: "" } });
  }, [isOpen]);
  const onValueChangeName = (value: string) => {
    dispatch({ type: "name", value });
  };
  const onClearName = () => {
    dispatch({ type: "name", value: "" });
  };
  const onValueChangeIntro = (value: string) => {
    dispatch({ type: "intro", value });
  };
  const onClearIntro = () => {
    dispatch({ type: "intro", value: "" });
  };
  const onPressConfirm = async () => {
    try {
      const resp: any = await addRole({
        ...state,
      });
      addToast({
        title: t("app.prompt.ok"),
        description: resp.msg,
        color: "success",
      });
      onClose();
    } catch (e) {
      addToast({
        title: (e as Error).cause as string,
        description: (e as Error).message,
        color: "danger",
      });
    }
  };
  return (
    <Modal isOpen={isOpen} placement='top-center' onOpenChange={onOpenChange}>
      <ModalContent>
        {(onClose) => (
          <>
            <ModalHeader className='flex flex-col gap-1'>
              {t("role.add.header")}
            </ModalHeader>
            <ModalBody>
              <Input
                endContent={
                  <UserIcon className='text-2xl text-default-400 pointer-events-none shrink-0' />
                }
                label={t("role.label.name")}
                isRequired
                type='text'
                placeholder={t("role.placeholder.name")}
                variant='bordered'
                onValueChange={onValueChangeName}
                onClear={onClearName}
              />
              <Input
                endContent={
                  <DocumentTextIcon className='text-2xl text-default-400 pointer-events-none shrink-0' />
                }
                label={t("role.label.intro")}
                isRequired
                placeholder={t("role.placeholder.intro")}
                variant='bordered'
                onValueChange={onValueChangeIntro}
                onClear={onClearIntro}
              />
            </ModalBody>
            <ModalFooter>
              <Button color='danger' variant='flat' onPress={onClose}>
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

const EditRole = ({
  role,
  isOpen,
  onClose,
  onOpenChange,
}: {
  role: TRole | null;
  isOpen: boolean;
  onClose: () => void;
  onOpenChange: () => void;
  isControlled?: boolean;
  getButtonProps?: (props?: any) => any;
  getDisclosureProps?: (props?: any) => any;
}) => {
  const { t } = useTranslation();
  const reducer = (state: any, action: any) => {
    switch (action.type) {
      case "name":
        return { ...state, name: action.value };
      case "intro":
        return { ...state, intro: action.value };
      default:
        return { ...state, ...action.value };
    }
  };
  const [state, dispatch] = React.useReducer(reducer, {
    xid: role?.xid ?? "",
    name: role?.name ?? "",
    intro: role?.intro ?? "",
  });
  React.useEffect(() => {
    dispatch({ type: "xid", value: role?.xid ?? "" });
    dispatch({ type: "name", value: role?.name ?? "" });
    dispatch({ type: "intro", value: role?.intro ?? "" });
  }, [role]);
  const onValueChangeName = (value: string) => {
    dispatch({ type: "name", value });
  };
  const onClearName = () => {
    dispatch({ type: "name", value: "" });
  };
  const onValueChangeIntro = (value: string) => {
    dispatch({ type: "intro", value });
  };
  const onClearIntro = () => {
    dispatch({ type: "intro", value: "" });
  };
  const onPressConfirm = async () => {
    try {
      if (!role) {
        return;
      }
      const resp: any = await editRole({
        ...state,
        xid: role.xid,
      });
      addToast({
        title: t("app.prompt.ok"),
        description: resp.msg,
        color: "success",
      });
      onClose();
    } catch (e) {
      addToast({
        title: (e as Error).cause as string,
        description: (e as Error).message,
        color: "danger",
      });
    }
  };
  React.useEffect(() => {
    dispatch({ value: { name: role?.name ?? "", intro: role?.intro ?? "" } });
  }, [role]);
  return (
    <Modal isOpen={isOpen} placement='top-center' onOpenChange={onOpenChange}>
      <ModalContent>
        {(onClose) => (
          <>
            <ModalHeader className='flex flex-col gap-1'>
              {t("role.edit.header")}
            </ModalHeader>
            <ModalBody>
              <Input
                endContent={
                  <UserIcon className='text-2xl text-default-400 pointer-events-none shrink-0' />
                }
                label={t("role.label.name")}
                isRequired
                type='text'
                placeholder={t("role.placeholder.name")}
                variant='bordered'
                onValueChange={onValueChangeName}
                onClear={onClearName}
                defaultValue={state.name}
                value={state.name}
              />
              <Input
                endContent={
                  <DocumentTextIcon className='text-2xl text-default-400 pointer-events-none shrink-0' />
                }
                label={t("role.label.intro")}
                isRequired
                type='text'
                placeholder={t("role.placeholder.intro")}
                variant='bordered'
                onValueChange={onValueChangeIntro}
                onClear={onClearIntro}
                defaultValue={state.intro}
                value={state.intro}
              />
            </ModalBody>
            <ModalFooter>
              <Button color='danger' variant='flat' onPress={onClose}>
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

const DeleteRole = ({
  role,
  isOpen,
  onClose,
  onOpenChange,
}: {
  role: TRole | null;
  isOpen: boolean;
  onClose: () => void;
  onOpenChange: () => void;
  isControlled?: boolean;
  getButtonProps?: (props?: any) => any;
  getDisclosureProps?: (props?: any) => any;
}) => {
  const { t } = useTranslation();
  const onPressConfirm = async () => {
    try {
      if (!role) {
        return;
      }
      const resp: any = await deleteRole(role.xid);
      addToast({
        title: t("app.prompt.ok"),
        description: resp.msg,
        color: "success",
      });
      onClose();
    } catch (e) {
      addToast({
        title: (e as Error).cause as string,
        description: (e as Error).message,
        color: "danger",
      });
    }
  };
  return (
    <Modal isOpen={isOpen} placement='top-center' onOpenChange={onOpenChange}>
      <ModalContent>
        {(onClose) => (
          <>
            <ModalHeader className='flex flex-col gap-1'>
              {t("role.edit.header")}
            </ModalHeader>
            <ModalBody>
              <p>{t("app.prompt.delete", { target: role?.name ?? "" })}</p>
            </ModalBody>
            <ModalFooter>
              <Button color='danger' variant='flat' onPress={onClose}>
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

export default function RolePage() {
  const searchParams = useSearchParams();
  const { t } = useTranslation();
  const columns = [
    {
      key: "xid",
      label: "role.label.id",
    },
    {
      key: "name",
      label: "role.label.name",
    },
    {
      key: "intro",
      label: "role.label.intro",
    },
    {
      key: "created_at",
      label: "role.label.createdAt",
    },
    {
      key: "updated_at",
      label: "role.label.updatedAt",
    },
    {
      key: "operation",
      label: "role.label.operation",
    },
  ];
  const reducer = (state: any, action: any) => {
    switch (action.type) {
      case "page":
        return { ...state, page: action.value };
      case "pageSize":
        return { ...state, pageSize: action.value };
      case "total":
        return { ...state, total: action.value };
      case "name":
        return { ...state, name: action.value };
      case "dxid":
        return { ...state, dxid: action.value };
      case "list":
        return { ...state, list: action.value };
      default:
        return { ...state, ...action.value };
    }
  };
  const [state, dispatch] = React.useReducer(reducer, {
    page: 1, // start from 1
    pageSize: 10,
    name: "",
    dxid: searchParams.get("dxid") ?? "",
    list: [],
    total: 0,
  });
  const onChangePagination = (page: number) => {
    dispatch({ type: "page", value: page });
    searchRole(page);
  };
  const [isLoading, setIsLoading] = React.useState(false);
  const searchRole = async (page: number) => {
    try {
      setIsLoading(true);
      const resp: any = await listRole({
        page,
        pageSize: state.pageSize,
        name: state.name,
        dxid: state.dxid,
      });
      setIsLoading(false);
      dispatch({ type: "list", value: resp.data.list ?? [] });
      dispatch({ type: "total", value: resp.data.total / state.pageSize + 1 });
      addToast({
        title: t("app.prompt.ok"),
        description: resp.msg,
        color: "success",
      });
    } catch (e) {
      setIsLoading(false);
      addToast({
        title: (e as Error).cause as string,
        description: (e as Error).message,
        color: "danger",
      });
    }
  };
  const token = useAuthStore((state: AuthStore) => state.token);
  React.useEffect(() => {
    if (!token) {
      // reset states
      dispatch({
        value: {
          page: 1, // start from 1
          pageSize: 10,
          name: "",
          dxid: searchParams.get("dxid") ?? "",
          list: [],
          total: 0,
        },
      });
      return;
    }
    searchRole(state.page);
  }, [token]);

  const onValueChangeName = (value: string) => {
    dispatch({ type: "name", value });
  };
  const onClearName = () => {
    dispatch({ type: "name", value: "" });
  };
  const onPressSearch = () => {
    searchRole(state.page);
  };

  // add role modal
  const {
    isOpen: isOpenAdd,
    onOpen: onOpenAdd,
    onOpenChange: onOpenChangeAdd,
    onClose: onCloseAdd,
  } = useDisclosure();
  const onCloseAddRole = () => {
    onCloseAdd();
    searchRole(state.page);
  };
  const [role, setRole] = React.useState<TRole | null>(null);
  // edit role modal
  const {
    isOpen: isOpenEdit,
    onOpen: onOpenEdit,
    onOpenChange: onOpenChangeEdit,
    onClose: onCloseEdit,
  } = useDisclosure();
  const onPressEdit = (role: TRole) => () => {
    setRole(role);
    onOpenEdit();
  };
  const onCloseEditRole = () => {
    onCloseEdit();
    searchRole(state.page);
  };
  // delete role modal
  const {
    isOpen: isOpenDelete,
    onOpen: onOpenDelete,
    onOpenChange: onOpenChangeDelete,
    onClose: onCloseDelete,
  } = useDisclosure();
  const onPressDelete = (role: TRole) => () => {
    setRole(role);
    onOpenDelete();
  };
  const onCloseDeleteRole = () => {
    onCloseDelete();
    searchRole(state.page);
  };
  return (
    <>
      {token && (
        <>
          <div className='flex flex-col items-center justify-center text-center'>
            <div className='w-[100%] flex flex-row flex-wrap text-start m-1'>
              {hasWidget(siteWidgets.btn_search_role) && (
                <>
                  <Input
                    className='w-[16%] m-1 h-10'
                    label={t("role.label.name")}
                    size='sm'
                    type='name'
                    variant='underlined'
                    onValueChange={onValueChangeName}
                    onClear={onClearName}
                  />
                  <Button
                    className='w-[80px] bg-blue-400 text-white m-1 h-10'
                    isIconOnly
                    startContent={
                      <FunnelIcon size={16} className='text-white mr-2' />
                    }
                    onPress={onPressSearch}
                  >
                    {t("app.btn.search")}
                  </Button>
                  {hasWidget(siteWidgets.btn_add_role) && (
                    <Button
                      className='w-[80px] bg-green-400 text-white m-1 h-10'
                      isIconOnly
                      startContent={
                        <PlusIcon size={16} className='text-white mr-2' />
                      }
                      onPress={onOpenAdd}
                    >
                      {t("app.btn.add")}
                    </Button>
                  )}
                </>
              )}
            </div>
            <Table
              className='min-w-[90vw] max-w-screen min-h-[700px] max-h-[1200px]'
              aria-label='role table'
            >
              <TableHeader columns={columns}>
                {(column) => (
                  <TableColumn key={column.key}>{t(column.label)}</TableColumn>
                )}
              </TableHeader>
              <TableBody
                items={state.list}
                emptyContent={t("app.prompt.tableNoContent")}
              >
                {(item: TRole) => (
                  <TableRow key={item.xid}>
                    {(columnKey) => {
                      if (columnKey != "operation") {
                        return (
                          <TableCell>{getKeyValue(item, columnKey)}</TableCell>
                        );
                      } else {
                        return (
                          <TableCell>
                            <Button
                              className='bg-transparent'
                              isIconOnly
                              startContent={
                                <PencilSquareIcon
                                  size={20}
                                  className='text-blue-400'
                                />
                              }
                              onPress={onPressEdit(item)}
                            />
                            <Button
                              className='bg-transparent'
                              isIconOnly
                              startContent={
                                <TrashIcon size={20} className='text-red-400' />
                              }
                              onPress={onPressDelete(item)}
                            />
                          </TableCell>
                        );
                      }
                    }}
                  </TableRow>
                )}
              </TableBody>
            </Table>

            {state.list.length > 0 && (
              <Pagination
                initialPage={state.page ?? 1}
                total={state.total}
                onChange={onChangePagination}
              />
            )}
          </div>
          <AddRole
            isOpen={isOpenAdd}
            onOpenChange={onOpenChangeAdd}
            onClose={onCloseAddRole}
          />
          <EditRole
            role={role}
            isOpen={isOpenEdit}
            onOpenChange={onOpenChangeEdit}
            onClose={onCloseEditRole}
          />
          <DeleteRole
            role={role}
            isOpen={isOpenDelete}
            onOpenChange={onOpenChangeDelete}
            onClose={onCloseDeleteRole}
          />
        </>
      )}
    </>
  );
}
