import type { ReactNode } from "react";

import { QueryProvider } from "@/components/query-provider";

import "./globals.css";

export const metadata = {
  title: "TraceDrop",
  description: "Developer observability for trace timelines",
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body>
        <QueryProvider>{children}</QueryProvider>
      </body>
    </html>
  );
}
