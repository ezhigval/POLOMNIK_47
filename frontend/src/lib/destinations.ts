export type Destination = {
  id: string;
  label: string;
  region: string;
  search: string;
  image: string;
};

export const departureCities = [
  "Санкт-Петербург",
  "Москва",
  "Казань",
  "Нижний Новгород",
] as const;

export const popularDestinations: Destination[] = [
  {
    id: "tikhvin",
    label: "Тихвинский путь",
    region: "Ленинградская область",
    search: "Тихвин",
    image: "https://images.unsplash.com/photo-1548013146-724f68d1ddac?auto=format&fit=crop&w=800&q=80",
  },
  {
    id: "optina",
    label: "Оптина пустынь",
    region: "Калужская область",
    search: "Оптина",
    image: "https://images.unsplash.com/photo-1548013146-724f68d1ddac?auto=format&fit=crop&w=800&q=80",
  },
  {
    id: "diveevo",
    label: "Дивеево",
    region: "Нижегородская область",
    search: "Дивеево",
    image: "https://images.unsplash.com/photo-1605647540924-852290f6b0d5?auto=format&fit=crop&w=800&q=80",
  },
  {
    id: "valaam",
    label: "Валаам",
    region: "Карелия",
    search: "Валаам",
    image: "https://images.unsplash.com/photo-1516026672322-bc52d61a55d5?auto=format&fit=crop&w=800&q=80",
  },
  {
    id: "solovki",
    label: "Соловки",
    region: "Архангельская область",
    search: "Солов",
    image: "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?auto=format&fit=crop&w=800&q=80",
  },
];

export function findDestination(id: string): Destination | undefined {
  return popularDestinations.find((item) => item.id === id);
}
