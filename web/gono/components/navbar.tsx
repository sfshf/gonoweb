"use client";
import {
  Navbar as HeroUINavbar,
  NavbarContent,
  NavbarMenu,
  NavbarBrand,
  NavbarItem,
  NavbarMenuItem,
} from "@heroui/navbar";
import { Button } from "@heroui/button";
import { Kbd } from "@heroui/kbd";
import { Link } from "@heroui/link";
import { Input } from "@heroui/input";
import { link as linkStyles } from "@heroui/theme";
import NextLink from "next/link";
import clsx from "clsx";
import React from "react";
import { hasMenu, siteConfig } from "@/config/site";
import { ThemeSwitch } from "@/components/theme-switch";
import { GlobeIcon } from "@/components/icons";
import {
  GithubIcon,
  UserCircleIcon,
  SearchIcon,
  Logo,
} from "@/components/icons";
import { useAuthStore } from "@/zustand/user";
import {
  Dropdown,
  DropdownTrigger,
  DropdownMenu,
  DropdownItem,
  addToast,
} from "@heroui/react";
import { signOut } from "@/api/user";
import { AuthStore, TResource } from "@/zustand/types";
import { User } from "@heroui/react";
import { useTranslation } from "react-i18next";
import i18next from "i18next";
import { usePathname } from "next/navigation";

export const Navbar = () => {
  const { t } = useTranslation();
  const pathname = usePathname();
  // 检查用户登录状态
  const token = useAuthStore((state: AuthStore) => state.token);
  const user = useAuthStore((state: AuthStore) => state.user);
  const menus = useAuthStore((state: AuthStore) => state.menus);

  const searchInput = (
    <Input
      aria-label='Search'
      classNames={{
        inputWrapper: "bg-default-100",
        input: "text-sm",
      }}
      endContent={
        <Kbd className='hidden lg:inline-block' keys={["command"]}>
          K
        </Kbd>
      }
      labelPlacement='outside'
      placeholder={t("app.label.search")}
      startContent={
        <SearchIcon className='text-base text-default-400 pointer-events-none flex-shrink-0' />
      }
      type='search'
    />
  );
  const clearAuth = useAuthStore((state: AuthStore) => state.clearAuth);
  const onPressSignOut = async () => {
    try {
      await signOut();
    } catch (e) {
    } finally {
      clearAuth();
      addToast({
        title: t("app.prompt.ok"),
        color: "success",
      });
    }
  };
  const [langs, setLangs] = React.useState<readonly string[]>(
    Array.isArray(i18next.options.supportedLngs)
      ? i18next.options.supportedLngs
      : [],
  );
  const onPressLang = (lang: string) => () => {
    i18next.changeLanguage(lang, (err, t) => {});
  };

  return (
    <HeroUINavbar maxWidth='xl' position='sticky'>
      <NavbarContent className='basis-1/5 sm:basis-full' justify='start'>
        <NavbarBrand as='li' className='gap-3 max-w-fit'>
          <NextLink className='flex justify-start items-center gap-1' href='/'>
            <Logo />
            <p className='font-bold text-inherit'>GONO</p>
          </NextLink>
        </NavbarBrand>
        <ul className='hidden lg:flex gap-4 justify-start ml-2'>
          {menus &&
            menus.map((item: TResource) => {
              const menu = hasMenu(item.identifier);
              if (!menu) {
                return <></>;
              }
              return (
                <NavbarItem key={menu.href}>
                  <NextLink
                    className={clsx(
                      pathname === menu.href
                        ? " text-blue-500 font-medium"
                        : "text-gray-700 hover:text-blue-500 ",
                    )}
                    color='foreground'
                    href={menu.href}
                  >
                    {t("app.menus." + menu.label)}
                  </NextLink>
                </NavbarItem>
              );
            })}
        </ul>
      </NavbarContent>

      <NavbarContent
        className='hidden sm:flex basis-1/5 sm:basis-full'
        justify='end'
      >
        <NavbarItem className='hidden sm:flex gap-2'>
          <Link isExternal aria-label='Github' href={siteConfig.links.github}>
            <GithubIcon className='text-default-500' />
          </Link>
          <ThemeSwitch />
          <Dropdown>
            <DropdownTrigger>
              <GlobeIcon
                className='cursor-pointer text-default-500 opacity-100 hover:opacity-80'
                size={22}
              />
            </DropdownTrigger>
            <DropdownMenu aria-label='Static Actions'>
              {langs.map((item) => {
                if (item != "cimode") {
                  return (
                    <DropdownItem key={item} onPress={onPressLang(item)}>
                      {t("app.lang." + item)}
                    </DropdownItem>
                  );
                } else {
                  return null;
                }
              })}
            </DropdownMenu>
          </Dropdown>
        </NavbarItem>
        <NavbarItem className='hidden lg:flex'>{searchInput}</NavbarItem>
        <NavbarItem className='hidden md:flex'>
          {!token && (
            <Button
              as={Link}
              className='text-sm font-normal text-default-600 bg-default-100'
              href={siteConfig.links.signIn}
              startContent={<UserCircleIcon className='text-blue' />}
              variant='flat'
            >
              {t("app.label.signIn")}
            </Button>
          )}
          {token && user && (
            <Dropdown>
              <DropdownTrigger>
                <User
                  avatarProps={{
                    src: user.avatar,
                  }}
                  description={user.email}
                  name={user.nick_name}
                />
              </DropdownTrigger>
              <DropdownMenu aria-label='Static Actions'>
                <DropdownItem
                  key='delete'
                  className='text-danger'
                  color='danger'
                  onPress={onPressSignOut}
                >
                  {t("app.label.signOut")}
                </DropdownItem>
              </DropdownMenu>
            </Dropdown>
          )}
        </NavbarItem>
      </NavbarContent>

      <NavbarMenu>
        {searchInput}
        <div className='mx-4 mt-2 flex flex-col gap-2'>
          {siteConfig.navMenuItems.map((item, index) => (
            <NavbarMenuItem key={`${item}-${index}`}>
              <Link
                color={
                  index === 2
                    ? "primary"
                    : index === siteConfig.navMenuItems.length - 1
                      ? "danger"
                      : "foreground"
                }
                href='#'
                size='lg'
              >
                {item.label}
              </Link>
            </NavbarMenuItem>
          ))}
        </div>
      </NavbarMenu>
    </HeroUINavbar>
  );
};
