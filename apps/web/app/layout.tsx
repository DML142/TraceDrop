import type { ReactNode } from "react";

export const metadata = {
  title: "TraceDrop",
  description: "Developer observability for trace timelines",
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
