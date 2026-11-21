import { siteConfig } from "@/lib/site-config";

import { russianPhone } from "@/lib/ru-contact";

export const contactPhone = russianPhone(siteConfig.contactPhone).href.slice(4);
export const contactPhoneDisplay = russianPhone(siteConfig.contactPhoneDisplay).label;
export const contactEmail = siteConfig.contactEmail;
