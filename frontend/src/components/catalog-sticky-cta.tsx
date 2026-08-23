import { contactPhone, contactPhoneDisplay } from "@/lib/contact";

export function CatalogStickyCTA() {
  return (
    <div className="fixed inset-x-0 bottom-0 z-30 border-t border-stone-200 bg-white/95 p-3 backdrop-blur-md lg:hidden">
      <div className="mx-auto flex max-w-6xl items-center gap-3 px-1">
        <a href={`tel:${contactPhone}`} className="btn-primary min-w-0 flex-1">
          {contactPhoneDisplay}
        </a>
        <a href="/support/chat" className="btn-secondary shrink-0">
          Чат
        </a>
      </div>
    </div>
  );
}
