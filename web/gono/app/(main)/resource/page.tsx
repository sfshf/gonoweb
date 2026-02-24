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
  Select,
  SelectItem,
  Image,
} from "@heroui/react";
import {
  listResource,
  addResource,
  editResource,
  deleteResource,
} from "@/api/resource";
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
import { AuthStore, TResource } from "@/zustand/types";
import { siteWidgets } from "@/config/site";
import { Link } from "@heroui/link";

const resourceTypes = ["menu", "widget", "api"];
const resourceType = (label: string): number => {
  switch (label) {
    case "menu":
      return 1;
    case "widget":
      return 2;
    case "api":
      return 3;
  }
  return -1;
};
const resourceTypeLabel = (typ: number): string => {
  switch (typ) {
    case 1:
      return "menu";
    case 2:
      return "widget";
    case 3:
      return "api";
  }
  return "";
};

const AddResource = ({
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
      case "type":
        return { ...state, type: action.value };
      case "identifier":
        return { ...state, identifier: action.value };
      case "name":
        return { ...state, name: action.value };
      case "intro":
        return { ...state, intro: action.value };
      case "icon":
        return { ...state, icon: action.value };
      default:
        return { ...state, ...action.value };
    }
  };
  const [state, dispatch] = React.useReducer(reducer, {
    type: -1,
    identifier: "",
    name: "",
    intro: "",
  });
  React.useEffect(() => {
    dispatch({ value: { type: -1, identifier: "", name: "", intro: "" } });
  }, [isOpen]);

  const onChangeType = (e: any) => {
    dispatch({ type: "type", value: resourceType(e.target.value) });
  };
  const onValueChangeIdentifier = (value: string) => {
    dispatch({ type: "identifier", value });
  };
  const onClearIdentifier = () => {
    dispatch({ type: "identifier", value: "" });
  };
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
  const onValueChangeIcon = (value: string) => {
    dispatch({ type: "icon", value });
  };
  const onClearIcon = () => {
    dispatch({ type: "icon", value: "" });
  };
  const onPressConfirm = async () => {
    try {
      const resp: any = await addResource({
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
              {t("resource.add.header")}
            </ModalHeader>
            <ModalBody>
              <Select
                label={t("resource.placeholder.type")}
                isClearable={true}
                onChange={onChangeType}
              >
                {resourceTypes.map((type) => (
                  <SelectItem key={type}>
                    {t("resource.type." + type)}
                  </SelectItem>
                ))}
              </Select>
              <Input
                endContent={
                  <UserIcon className='text-2xl text-default-400 pointer-events-none shrink-0' />
                }
                label={t("resource.label.id")}
                isRequired
                type='text'
                placeholder={t("resource.placeholder.id")}
                variant='bordered'
                onValueChange={onValueChangeIdentifier}
                onClear={onClearIdentifier}
              />
              <Input
                endContent={
                  <UserIcon className='text-2xl text-default-400 pointer-events-none shrink-0' />
                }
                label={t("resource.label.name")}
                isRequired
                type='text'
                placeholder={t("resource.placeholder.name")}
                variant='bordered'
                onValueChange={onValueChangeName}
                onClear={onClearName}
              />
              <Input
                endContent={
                  <DocumentTextIcon className='text-2xl text-default-400 pointer-events-none shrink-0' />
                }
                label={t("resource.label.intro")}
                isRequired
                placeholder={t("resource.placeholder.intro")}
                variant='bordered'
                onValueChange={onValueChangeIntro}
                onClear={onClearIntro}
              />
              <Input
                endContent={
                  <DocumentTextIcon className='text-2xl text-default-400 pointer-events-none shrink-0' />
                }
                label={t("resource.label.icon")}
                isRequired
                placeholder={t("resource.placeholder.icon")}
                variant='bordered'
                onValueChange={onValueChangeIcon}
                onClear={onClearIcon}
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

const EditResource = ({
  resource,
  isOpen,
  onClose,
  onOpenChange,
}: {
  resource: TResource | null;
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
      case "type":
        return { ...state, type: action.value };
      case "identifier":
        return { ...state, identifier: action.value };
      case "name":
        return { ...state, name: action.value };
      case "intro":
        return { ...state, intro: action.value };
      case "icon":
        return { ...state, icon: action.value };
      default:
        return { ...state, ...action.value };
    }
  };
  const [state, dispatch] = React.useReducer(reducer, {
    type: resource?.type ?? -1,
    identifier: resource?.identifier ?? "",
    name: resource?.name ?? "",
    intro: resource?.intro ?? "",
    icon: resource?.icon ?? "",
  });
  React.useEffect(() => {
    dispatch({
      value: {
        xid: resource?.identifier ?? "",
        name: resource?.name ?? "",
        intro: resource?.intro ?? "",
      },
    });
  }, [resource]);
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
  const onValueChangeIcon = (value: string) => {
    dispatch({ type: "icon", value });
  };
  const onClearIcon = () => {
    dispatch({ type: "icon", value: "" });
  };
  const onPressConfirm = async () => {
    try {
      if (!resource) {
        return;
      }
      const resp: any = await editResource({
        ...state,
        id: resource.id,
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
    dispatch({
      value: {
        name: resource?.name ?? "",
        type: resource?.type,
        identifier: resource?.identifier ?? "",
        intro: resource?.intro ?? "",
        icon: resource?.icon ?? "",
      },
    });
  }, [resource]);
  return (
    <Modal isOpen={isOpen} placement='top-center' onOpenChange={onOpenChange}>
      <ModalContent>
        {(onClose) => (
          <>
            <ModalHeader className='flex flex-col gap-1'>
              {t("resource.edit.header")}
            </ModalHeader>
            <ModalBody>
              <Select
                disabled
                label={t("resource.placeholder.type")}
                selectedKeys={[resourceTypeLabel(state.type)]}
              >
                {resourceTypes.map((type) => (
                  <SelectItem key={type}>
                    {t("resource.type." + type)}
                  </SelectItem>
                ))}
              </Select>
              <Input
                disabled
                value={state.identifier}
                endContent={
                  <UserIcon className='text-2xl text-default-400 pointer-events-none shrink-0' />
                }
                label={t("resource.label.id")}
                isRequired
                type='text'
                placeholder={t("resource.placeholder.id")}
                variant='bordered'
              />
              <Input
                endContent={
                  <UserIcon className='text-2xl text-default-400 pointer-events-none shrink-0' />
                }
                label={t("resource.label.name")}
                isRequired
                type='text'
                placeholder={t("resource.placeholder.name")}
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
                label={t("resource.label.intro")}
                isRequired
                type='text'
                placeholder={t("resource.placeholder.intro")}
                variant='bordered'
                onValueChange={onValueChangeIntro}
                onClear={onClearIntro}
                defaultValue={state.intro}
                value={state.intro}
              />
              <Input
                endContent={
                  <DocumentTextIcon className='text-2xl text-default-400 pointer-events-none shrink-0' />
                }
                label={t("resource.label.icon")}
                isRequired
                type='text'
                placeholder={t("resource.placeholder.icon")}
                variant='bordered'
                onValueChange={onValueChangeIcon}
                onClear={onClearIcon}
                defaultValue={state.icon}
                value={state.icon}
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

const DeleteResource = ({
  resource,
  isOpen,
  onClose,
  onOpenChange,
}: {
  resource: TResource | null;
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
      if (!resource) {
        return;
      }
      const resp: any = await deleteResource(resource.id);
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
              {t("resource.edit.header")}
            </ModalHeader>
            <ModalBody>
              <p>
                {t("app.prompt.delete", {
                  target: resource?.identifier ?? "",
                })}
              </p>
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

export default function ResourcePage() {
  const { t } = useTranslation();
  const columns = [
    {
      key: "identifier",
      label: "resource.label.id",
    },
    {
      key: "type",
      label: "resource.label.type",
    },
    {
      key: "name",
      label: "resource.label.name",
    },
    {
      key: "icon",
      label: "resource.label.icon",
    },
    {
      key: "intro",
      label: "resource.label.intro",
    },
    {
      key: "created_at",
      label: "resource.label.createdAt",
    },
    {
      key: "updated_at",
      label: "resource.label.updatedAt",
    },
    {
      key: "operation",
      label: "resource.label.operation",
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
      case "type":
        return { ...state, type: action.value };
      case "name":
        return { ...state, name: action.value };
      case "identifier":
        return { ...state, identifier: action.value };
      case "list":
        return { ...state, list: action.value };
      default:
        return { ...state, ...action.value };
    }
  };
  const [state, dispatch] = React.useReducer(reducer, {
    page: 1, // start from 1
    pageSize: 10,
    type: -1,
    identifier: "",
    name: "",
    list: [],
    total: 0,
  });
  const onChangePagination = (page: number) => {
    dispatch({ type: "page", value: page });
    searchResource(page);
  };
  const [isLoading, setIsLoading] = React.useState(false);
  const searchResource = async (page: number) => {
    try {
      setIsLoading(true);
      const resp: any = await listResource({
        page,
        pageSize: state.pageSize,
        type: state.type,
        identifier: state.identifier,
        name: state.name,
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
  const token = useAuthStore((state: AuthStore) => state.token);
  React.useEffect(() => {
    if (!token) {
      // reset states
      dispatch({ type: "list", value: [] });
      dispatch({ type: "total", value: 0 });
      return;
    }
    searchResource(state.page);
  }, [token]);
  const onChangeType = (e: any) => {
    dispatch({ type: "type", value: resourceType(e.target.value) });
  };
  const onValueChangeName = (value: string) => {
    dispatch({ type: "name", value });
  };
  const onClearName = () => {
    dispatch({ type: "name", value: "" });
  };
  const onValueChangeIdentifier = (value: string) => {
    dispatch({ type: "identifier", value });
  };
  const onClearIdentifier = () => {
    dispatch({ type: "identifier", value: "" });
  };
  const onPressSearch = () => {
    searchResource(state.page);
  };

  // add resource modal
  const {
    isOpen: isOpenAdd,
    onOpen: onOpenAdd,
    onOpenChange: onOpenChangeAdd,
    onClose: onCloseAdd,
  } = useDisclosure();
  const onCloseAddResource = () => {
    onCloseAdd();
    searchResource(state.page);
  };
  // edit resource modal
  const {
    isOpen: isOpenEdit,
    onOpen: onOpenEdit,
    onOpenChange: onOpenChangeEdit,
    onClose: onCloseEdit,
  } = useDisclosure();
  const [resource, setResource] = React.useState<TResource | null>(null);
  const onPressEdit = (resource: TResource) => () => {
    setResource(resource);
    onOpenEdit();
  };
  const onCloseEditResource = () => {
    onCloseEdit();
    searchResource(state.page);
  };
  // delete resource modal
  const {
    isOpen: isOpenDelete,
    onOpen: onOpenDelete,
    onOpenChange: onOpenChangeDelete,
    onClose: onCloseDelete,
  } = useDisclosure();
  const onPressDelete = (resource: TResource) => () => {
    setResource(resource);
    onOpenDelete();
  };
  const onCloseDeleteResource = () => {
    onCloseDelete();
    searchResource(state.page);
  };
  return (
    <>
      {token && (
        <>
          <div className='flex flex-col items-center justify-center text-center'>
            <div className='w-[100%] flex flex-row flex-wrap text-start m-1'>
              {hasWidget(siteWidgets.btn_search_resource) && (
                <>
                  <Select
                    className='max-w-xs'
                    label={t("resource.placeholder.type")}
                    size='sm'
                    isClearable={true}
                    onChange={onChangeType}
                  >
                    {resourceTypes.map((type) => (
                      <SelectItem key={type}>
                        {t("resource.type." + type)}
                      </SelectItem>
                    ))}
                  </Select>
                  <Input
                    className='w-[16%] m-1 h-10'
                    label={t("resource.label.name")}
                    size='sm'
                    type='name'
                    variant='underlined'
                    onValueChange={onValueChangeName}
                    onClear={onClearName}
                  />
                  <Input
                    className='w-[16%] m-1 h-10'
                    label={t("resource.label.id")}
                    size='sm'
                    type='name'
                    variant='underlined'
                    onValueChange={onValueChangeIdentifier}
                    onClear={onClearIdentifier}
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
                  {hasWidget(siteWidgets.btn_add_resource) && (
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
              aria-label='resource table'
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
                {(item: TResource) => (
                  <TableRow key={item.identifier}>
                    {(columnKey) => {
                      if (columnKey != "operation") {
                        if (columnKey == "type") {
                          return (
                            <TableCell>
                              {t(
                                "resource.type." +
                                  resourceTypeLabel(
                                    getKeyValue(item, columnKey),
                                  ),
                              )}
                            </TableCell>
                          );
                        } else if (columnKey == "icon") {
                          return (
                            <TableCell>
                              <Link
                                isExternal
                                href={getKeyValue(item, columnKey)}
                              >
                                <Image
                                  isZoomed
                                  width={20}
                                  height={20}
                                  alt='HeroUI Fruit Image with Zoom'
                                  src={getKeyValue(item, columnKey)}
                                />
                              </Link>
                            </TableCell>
                          );
                        } else {
                          return (
                            <TableCell>
                              {getKeyValue(item, columnKey)}
                            </TableCell>
                          );
                        }
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
          <AddResource
            isOpen={isOpenAdd}
            onOpenChange={onOpenChangeAdd}
            onClose={onCloseAddResource}
          />
          <EditResource
            resource={resource}
            isOpen={isOpenEdit}
            onOpenChange={onOpenChangeEdit}
            onClose={onCloseEditResource}
          />
          <DeleteResource
            resource={resource}
            isOpen={isOpenDelete}
            onOpenChange={onOpenChangeDelete}
            onClose={onCloseDeleteResource}
          />
        </>
      )}
    </>
  );
}
