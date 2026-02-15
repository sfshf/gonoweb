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
import { listUser, addUser, editUser, deleteUser } from "@/api/user";
import { Pagination } from "@heroui/pagination";
import { useTranslation } from "react-i18next";
import {
  FunnelIcon,
  PencilSquareIcon,
  PlusIcon,
  MailIcon,
  TrashIcon,
  UserIcon,
} from "@/components/icons";
import { hasWidget, useAuthStore } from "@/zustand/user";
import { AuthStore, TUser } from "@/zustand/types";
import { siteWidgets } from "@/config/site";

const AddUser = ({
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
      case "email":
        return { ...state, email: action.value };
      case "nickname":
        return { ...state, nickname: action.value };
      default:
        return { ...action.value };
    }
  };
  const [state, dispatch] = React.useReducer(reducer, {
    email: "",
    nickname: "",
  });
  React.useEffect(() => {
    dispatch({ value: { email: "", nickname: "" } });
  }, [isOpen]);
  const onValueChangeEmail = (value: string) => {
    dispatch({ type: "email", value });
  };
  const onClearEmail = () => {
    dispatch({ type: "email", value: "" });
  };
  const onValueChangeNickname = (value: string) => {
    dispatch({ type: "nickname", value });
  };
  const onClearNickname = () => {
    dispatch({ type: "nickname", value: "" });
  };
  const validateEmail = (value: string) =>
    /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value);
  const onPressConfirm = async () => {
    try {
      if (!validateEmail(state.email)) {
        return;
      }
      const resp: any = await addUser({
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
    <>
      <Modal isOpen={isOpen} placement='top-center' onOpenChange={onOpenChange}>
        <ModalContent>
          {(onClose) => (
            <>
              <ModalHeader className='flex flex-col gap-1'>
                {t("user.add.header")}
              </ModalHeader>
              <ModalBody>
                <Input
                  endContent={
                    <MailIcon className='text-2xl text-default-400 pointer-events-none shrink-0' />
                  }
                  label={t("user.label.email")}
                  isRequired
                  type='email'
                  placeholder={t("user.placeholder.email")}
                  variant='bordered'
                  onValueChange={onValueChangeEmail}
                  onClear={onClearEmail}
                />
                <Input
                  endContent={
                    <UserIcon className='text-2xl text-default-400 pointer-events-none shrink-0' />
                  }
                  label={t("user.label.nickname")}
                  isRequired
                  placeholder={t("user.placeholder.nickname")}
                  variant='bordered'
                  onValueChange={onValueChangeNickname}
                  onClear={onClearNickname}
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
    </>
  );
};

const EditUser = ({
  user,
  isOpen,
  onClose,
  onOpenChange,
}: {
  user: TUser | null;
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
      case "email":
        return { ...state, email: action.value };
      case "nickname":
        return { ...state, nickname: action.value };
      default:
        return { ...action.value };
    }
  };
  const [state, dispatch] = React.useReducer(reducer, {
    xid: user?.xid ?? "",
    email: user?.email ?? "",
    nickname: user?.nick_name ?? "",
  });
  React.useEffect(() => {
    dispatch({ type: "xid", value: user?.xid ?? "" });
    dispatch({ type: "email", value: user?.email ?? "" });
    dispatch({ type: "nickname", value: user?.nick_name ?? "" });
  }, [user]);
  const onValueChangeEmail = (value: string) => {
    dispatch({ type: "email", value });
  };
  const onClearEmail = () => {
    dispatch({ type: "email", value: "" });
  };
  const onValueChangeNickname = (value: string) => {
    dispatch({ type: "nickname", value });
  };
  const onClearNickname = () => {
    dispatch({ type: "nickname", value: "" });
  };
  const validateEmail = (value: string) =>
    /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value);
  const onPressConfirm = async () => {
    try {
      if (!user) {
        return;
      }
      if (!validateEmail(state.email)) {
        return;
      }
      const resp: any = await editUser({
        ...state,
        xid: user.xid,
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
    <>
      <Modal isOpen={isOpen} placement='top-center' onOpenChange={onOpenChange}>
        <ModalContent>
          {(onClose) => (
            <>
              <ModalHeader className='flex flex-col gap-1'>
                {t("user.edit.header")}
              </ModalHeader>
              <ModalBody>
                <Input
                  endContent={
                    <MailIcon className='text-2xl text-default-400 pointer-events-none shrink-0' />
                  }
                  label={t("user.label.email")}
                  isRequired
                  type='email'
                  placeholder={t("user.placeholder.email")}
                  variant='bordered'
                  onValueChange={onValueChangeEmail}
                  onClear={onClearEmail}
                  defaultValue={state.email}
                  value={state.email}
                />
                <Input
                  endContent={
                    <UserIcon className='text-2xl text-default-400 pointer-events-none shrink-0' />
                  }
                  label={t("user.label.nickname")}
                  isRequired
                  placeholder={t("user.placeholder.nickname")}
                  variant='bordered'
                  onValueChange={onValueChangeNickname}
                  onClear={onClearNickname}
                  defaultValue={state.nickname}
                  value={state.nickname}
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
    </>
  );
};

const DeleteUser = ({
  user,
  isOpen,
  onClose,
  onOpenChange,
}: {
  user: TUser | null;
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
      if (!user) {
        return;
      }
      const resp: any = await deleteUser(user.xid);
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
    <>
      <Modal isOpen={isOpen} placement='top-center' onOpenChange={onOpenChange}>
        <ModalContent>
          {(onClose) => (
            <>
              <ModalHeader className='flex flex-col gap-1'>
                {t("user.edit.header")}
              </ModalHeader>
              <ModalBody>
                <p>{t("app.prompt.delete", { target: user?.email ?? "" })}</p>
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
    </>
  );
};

export default function UserPage() {
  const { t } = useTranslation();
  const columns = [
    {
      key: "xid",
      label: "user.label.id",
    },
    {
      key: "email",
      label: "user.label.email",
    },
    {
      key: "nick_name",
      label: "user.label.nickname",
    },
    {
      key: "avatar",
      label: "user.label.avatar",
    },
    {
      key: "real_name",
      label: "user.label.realname",
    },
    {
      key: "password",
      label: "user.label.password",
    },

    {
      key: "created_at",
      label: "user.label.createdAt",
    },
    {
      key: "updated_at",
      label: "user.label.updatedAt",
    },
    {
      key: "operation",
      label: "user.label.operation",
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
      case "email":
        return { ...state, email: action.value };
      case "nickname":
        return { ...state, nickname: action.value };
      case "realname":
        return { ...state, realname: action.value };
      case "list":
        return { ...state, list: action.value };
      default:
        return { ...action.value };
    }
  };
  const [state, dispatch] = React.useReducer(reducer, {
    page: 1, // start from 1
    pageSize: 10,
    email: "",
    nickname: "",
    realname: "",
    list: [],
    total: 0,
  });
  const onChangePagination = (page: number) => {
    dispatch({ type: "page", value: page });
    searchUser(page);
  };
  const [isLoading, setIsLoading] = React.useState(false);
  const searchUser = async (page: number) => {
    try {
      setIsLoading(true);
      const resp: any = await listUser({
        page,
        pageSize: state.pageSize,
        email: state.email,
        nickname: state.nickname,
        realname: state.realname,
      });
      setIsLoading(false);
      dispatch({ type: "list", value: resp.data.list });
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
  const isRoot = (email: string): boolean => {
    if (email === process.env.NEXT_PUBLIC_ROOT_EMAIL) {
      return true;
    }
    return false;
  };
  const token = useAuthStore((state: AuthStore) => state.token);
  React.useEffect(() => {
    if (!token) {
      // reset states
      dispatch({ type: "list", value: [] });
      dispatch({ type: "total", value: 0 });
      return;
    }
    searchUser(state.page);
  }, [token]);

  const onValueChangeEmail = (value: string) => {
    dispatch({ type: "email", value });
  };
  const onClearEmail = () => {
    dispatch({ type: "email", value: "" });
  };
  const onValueChangeNickname = (value: string) => {
    dispatch({ type: "nickname", value });
  };
  const onClearNickname = () => {
    dispatch({ type: "nickname", value: "" });
  };
  const onValueChangeRealname = (value: string) => {
    dispatch({ type: "realname", value });
  };
  const onClearRealname = () => {
    dispatch({ type: "realname", value: "" });
  };
  const onPressSearch = () => {
    searchUser(state.page);
  };

  // add user modal
  const {
    isOpen: isOpenAdd,
    onOpen: onOpenAdd,
    onOpenChange: onOpenChangeAdd,
    onClose: onCloseAdd,
  } = useDisclosure();
  const onCloseAddUser = () => {
    onCloseAdd();
    searchUser(state.page);
  };
  // edit user modal
  const {
    isOpen: isOpenEdit,
    onOpen: onOpenEdit,
    onOpenChange: onOpenChangeEdit,
    onClose: onCloseEdit,
  } = useDisclosure();
  const [user, setUser] = React.useState<TUser | null>(null);
  const onPressEdit = (user: TUser) => () => {
    setUser(user);
    onOpenEdit();
  };
  const onCloseEditUser = () => {
    onCloseEdit();
    searchUser(state.page);
  };
  // delete user modal
  const {
    isOpen: isOpenDelete,
    onOpen: onOpenDelete,
    onOpenChange: onOpenChangeDelete,
    onClose: onCloseDelete,
  } = useDisclosure();
  const onPressDelete = (user: TUser) => () => {
    setUser(user);
    onOpenDelete();
  };
  const onCloseDeleteUser = () => {
    onCloseDelete();
    searchUser(state.page);
  };
  return (
    <>
      {token && (
        <>
          <div className='flex flex-col items-center justify-center text-center'>
            <div className='w-[100%] flex flex-row flex-wrap text-start m-1'>
              {hasWidget(siteWidgets.btn_search_user) && (
                <>
                  <Input
                    className='w-[16%] m-1 h-10'
                    label={t("user.label.email")}
                    size='sm'
                    type='email'
                    variant='underlined'
                    onValueChange={onValueChangeEmail}
                    onClear={onClearEmail}
                  />
                  <Input
                    className='w-[16%] m-1 h-10'
                    label={t("user.label.nickname")}
                    size='sm'
                    type='text'
                    variant='underlined'
                    onValueChange={onValueChangeNickname}
                    onClear={onClearNickname}
                  />
                  <Input
                    className='w-[16%] m-1 h-10'
                    label={t("user.label.realname")}
                    size='sm'
                    type='text'
                    variant='underlined'
                    onValueChange={onValueChangeRealname}
                    onClear={onClearRealname}
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
                  {hasWidget(siteWidgets.btn_add_user) && (
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
              aria-label='user table'
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
                {(item: TUser) => (
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
                                  className={
                                    !isRoot(item.email)
                                      ? "text-blue-400"
                                      : "text-gray-400"
                                  }
                                />
                              }
                              onPress={onPressEdit(item)}
                            />
                            <Button
                              className='bg-transparent'
                              isIconOnly
                              startContent={
                                <TrashIcon
                                  size={20}
                                  className={
                                    !isRoot(item.email)
                                      ? "text-red-400"
                                      : "text-gray-400"
                                  }
                                />
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
          <AddUser
            isOpen={isOpenAdd}
            onOpenChange={onOpenChangeAdd}
            onClose={onCloseAddUser}
          />
          <EditUser
            user={user}
            isOpen={isOpenEdit}
            onOpenChange={onOpenChangeEdit}
            onClose={onCloseEditUser}
          />
          <DeleteUser
            user={user}
            isOpen={isOpenDelete}
            onOpenChange={onOpenChangeDelete}
            onClose={onCloseDeleteUser}
          />
        </>
      )}
    </>
  );
}
