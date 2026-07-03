import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "OLT Migrate",
  description: "Ferramenta de migração de OLT",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="pt-BR" className="bg-[#f9fafb] text-gray-900" suppressHydrationWarning>
      <body className="min-h-screen bg-[#f9fafb] text-gray-900 font-sans flex flex-col antialiased selection:bg-blue-100 selection:text-blue-900" suppressHydrationWarning>
        
        {/* Simple Navbar */}
        <header className="bg-white border-b border-gray-200">
          <div className="max-w-6xl mx-auto px-6 h-16 flex items-center justify-between">
            <div className="flex items-center gap-2">
              <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center text-white">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" /></svg>
              </div>
              <span className="font-semibold tracking-tight text-lg">OLT Migrate</span>
            </div>
          </div>
        </header>

        {/* Main Workspace */}
        <main className="flex-1 max-w-6xl w-full mx-auto p-6 md:p-8">
          {children}
        </main>
        
      </body>
    </html>
  );
}
