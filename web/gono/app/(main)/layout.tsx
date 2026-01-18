import { Metadata, Viewport } from "next";
import { Link } from "@heroui/link";
import clsx from "clsx";

import { siteConfig } from "@/config/site";
import { fontSans } from "@/config/fonts";
import { Navbar } from "@/components/navbar";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className='relative flex flex-col h-screen'>
      <Navbar />
      <main className='container mx-auto max-w-7xl pt-16 px-6 flex-grow'>
        {children}
      </main>
    </div>
  );
}
