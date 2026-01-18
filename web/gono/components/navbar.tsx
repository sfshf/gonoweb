"use client";
import {
  Navbar as HeroUINavbar,
  NavbarContent,
  NavbarMenu,
  NavbarMenuToggle,
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
import {
  GithubIcon,
  UserCircleIcon,
  SearchIcon,
  Logo,
} from "@/components/icons";
import { useAuthStore } from "@/zustand/user";
import { useRouter } from "next/navigation";
import {
  Dropdown,
  DropdownTrigger,
  DropdownMenu,
  DropdownItem,
  addToast,
} from "@heroui/react";
import { signOut } from "@/api/user";
import { AuthStore, TMenuWidget } from "@/zustand/types";
import { User } from "@heroui/react";

export const Navbar = () => {
  const router = useRouter();
  // 检查用户登录状态
  const token = useAuthStore((state: AuthStore) => state.token);
  const user = useAuthStore((state: AuthStore) => state.user);
  const menus = useAuthStore((state: AuthStore) => state.menus);
  React.useEffect(() => {
    // 未登录，则跳转登录页
    if (!token) {
      router.push("/sign-in");
      return;
    }
  }, [token]);

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
      placeholder='Search...'
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
        title: "OK",
        color: "success",
      });
    }
  };

  return (
    <HeroUINavbar maxWidth='xl' position='sticky'>
      <NavbarContent className='basis-1/5 sm:basis-full' justify='start'>
        <NavbarBrand as='li' className='gap-3 max-w-fit'>
          <NextLink className='flex justify-start items-center gap-1' href='/'>
            <Logo />
            <p className='font-bold text-inherit'>ACME</p>
          </NextLink>
        </NavbarBrand>
        <ul className='hidden lg:flex gap-4 justify-start ml-2'>
          {menus &&
            menus.map((item: TMenuWidget) => {
              const menu = hasMenu(item.identifier);
              if (!menu) {
                return <></>;
              }
              return (
                <NavbarItem key={menu.href}>
                  <NextLink
                    className={clsx(
                      linkStyles({ color: "foreground" }),
                      "data-[active=true]:text-primary data-[active=true]:font-medium"
                    )}
                    color='foreground'
                    href={menu.href}
                  >
                    {menu.label}
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
              登录
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
                  登出
                </DropdownItem>
              </DropdownMenu>
            </Dropdown>
          )}
        </NavbarItem>
      </NavbarContent>

      <NavbarContent className='sm:hidden basis-1 pl-4' justify='end'>
        <Link isExternal aria-label='Github' href={siteConfig.links.github}>
          <GithubIcon className='text-default-500' />
        </Link>
        <ThemeSwitch />
        <NavbarMenuToggle />
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
