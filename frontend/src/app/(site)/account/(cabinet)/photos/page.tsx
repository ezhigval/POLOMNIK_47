import type { Metadata } from "next";
import { redirect } from "next/navigation";
import { UserPhotoUploadForm, type UserPhoto } from "@/components/account/user-photo-upload-form";
import { getApiBaseUrl } from "@/lib/api/base-url";
import { getAuthToken } from "@/lib/auth/session";

export const metadata: Metadata = {
  title: "Фотографии",
};

async function fetchMyPhotos(token: string): Promise<UserPhoto[]> {
  const response = await fetch(`${getApiBaseUrl()}/me/photos`, {
    headers: { Authorization: `Bearer ${token}` },
    cache: "no-store",
  });
  if (!response.ok) {
    return [];
  }
  const body = await response.json().catch(() => null);
  return Array.isArray(body?.data) ? body.data : [];
}

export default async function AccountPhotosPage() {
  const token = await getAuthToken();
  if (!token) {
    redirect("/account/login?returnUrl=%2Faccount%2Fphotos");
  }

  const photos = await fetchMyPhotos(token);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="font-display text-3xl font-semibold text-stone-900">Фотографии</h1>
        <p className="mt-2 text-sm text-stone-600">
          Загрузка фотографий с отдельным согласием на обработку и, при желании, на распространение.
        </p>
      </div>
      <UserPhotoUploadForm photos={photos} />
    </div>
  );
}
