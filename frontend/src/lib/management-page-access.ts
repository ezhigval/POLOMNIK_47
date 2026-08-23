import "server-only";

import { getManagementSession } from "@/lib/api/management";
import { sessionHasAnyPermission } from "@/lib/management-access";

export async function canAccessManagementPage(anyOf: string[]): Promise<boolean> {
  try {
    const session = await getManagementSession();
    return sessionHasAnyPermission(session, anyOf);
  } catch {
    return false;
  }
}
