"use client";
import {
  Card,
  CardHeader,
  CardBody,
  CardFooter,
  Form,
  Input,
  Button,
  addToast,
} from "@heroui/react";
import React from "react";
import { EyeFilledIcon, EyeSlashFilledIcon } from "@/components/icons";
import { useTranslation } from "react-i18next";
import { signIn } from "@/api/user";
import { useAuthStore } from "@/zustand/user";
import { useRouter } from "next/navigation";
import { AuthStore } from "@/zustand/types";

type SignInForm = {
  account: string;
  password: string;
};

export default function SignInPage() {
  const { t } = useTranslation();
  const router = useRouter();
  const [signInForm, setSignInForm] = React.useState<SignInForm>({
    account: "",
    password: "",
  });
  const [passwordVisible, setPasswordVisible] = React.useState(false);
  const toggleVisibility = () => setPasswordVisible(!passwordVisible);
  const onChangeAccount = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSignInForm({
      ...signInForm,
      account: e.target.value,
    });
  };
  const onClearAccount = () => {
    setSignInForm({
      ...signInForm,
      account: "",
    });
  };
  const onChangePassword = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSignInForm({
      ...signInForm,
      password: e.target.value,
    });
  };
  // 存储用户信息 -- 当前的角色和域租户，能访问到的菜单和控件
  const setAuth = useAuthStore((state: AuthStore) => state.setAuth);
  const onSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    try {
      e.preventDefault();
      const resp: any = await signIn({ ...signInForm });
      setAuth(
        resp.data.token,
        resp.data.user,
        resp.data.domain,
        resp.data.role,
        resp.data.menus,
        resp.data.widgets,
      );
      addToast({
        title: t("app.prompt.ok"),
        description: resp.msg,
        color: "success",
      });
      router.back();
    } catch (e) {
      addToast({
        title: (e as Error).cause as string,
        description: (e as Error).message,
        color: "danger",
      });
    }
  };
  const token = useAuthStore((state: AuthStore) => state.token);
  React.useEffect(() => {
    // 已登录，则回跳
    if (token) {
      router.push("/");
    }
  }, [token]);

  return (
    <Card className='min-w-[450px] max-w-[800px] min-h-[700px] max-h-[1200px] mt-14'>
      <CardHeader className='flex gap-3 justify-center'>
        <div className='flex flex-col'>
          <p className='text-3xl'>Gono Web</p>
        </div>
      </CardHeader>
      <CardBody>
        <Form className='w-full flex flex-col mt-14 gap-14' onSubmit={onSubmit}>
          <Input
            isRequired
            errorMessage='Please enter a valid nickname/email'
            label={t("signIn.label.account")}
            labelPlacement='outside'
            name='nickname/email'
            placeholder={t("signIn.placeholder.nickname/email")}
            type='text'
            onChange={onChangeAccount}
            onClear={onClearAccount}
          />
          <Input
            isRequired
            errorMessage='Please enter password'
            label={t("signIn.label.Password")}
            labelPlacement='outside'
            name='password'
            placeholder={t("signIn.placeholder.password")}
            type={passwordVisible ? "text" : "password"}
            endContent={
              <button
                aria-label='toggle password visibility'
                className='focus:outline-solid outline-transparent'
                type='button'
                onClick={toggleVisibility}
              >
                {passwordVisible ? (
                  <EyeSlashFilledIcon className='text-2xl text-default-400 pointer-events-none' />
                ) : (
                  <EyeFilledIcon className='text-2xl text-default-400 pointer-events-none' />
                )}
              </button>
            }
            onChange={onChangePassword}
          />
          <Button
            className='w-full'
            color={
              signInForm.account == "" || signInForm.password == ""
                ? "default"
                : "primary"
            }
            type='submit'
            disabled={signInForm.account == "" || signInForm.password == ""}
          >
            {t("signIn.btn.1")}
          </Button>
        </Form>
      </CardBody>
      <CardFooter></CardFooter>
    </Card>
  );
}
