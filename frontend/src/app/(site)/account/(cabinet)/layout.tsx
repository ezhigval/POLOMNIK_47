import { AccountNav } from "@/components/account-nav";

export default function AccountCabinetLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="mx-auto max-w-6xl space-y-6 px-4 py-8 sm:py-12">
      <AccountNav />
      {children}
    </div>
  );
}
