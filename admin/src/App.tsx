import { Authenticated, GitHubBanner, Refine } from "@refinedev/core";
import { RefineKbar, RefineKbarProvider } from "@refinedev/kbar";

import {
  ErrorComponent,
  notificationProvider,
  ThemedLayoutV2,
  ThemedSiderV2,
  ThemedTitleV2,
} from "@refinedev/antd";
import "@refinedev/antd/dist/reset.css";

import routerBindings, {
  CatchAllNavigate,
  NavigateToResource,
  UnsavedChangesNotifier,
} from "@refinedev/react-router-v6";
import dataProvider from "@refinedev/simple-rest";
import { useTranslation } from "react-i18next";
import { BrowserRouter, Outlet, Route, Routes } from "react-router-dom";
import { authProvider, TOKEN_KEY, USER_KEY } from "./authProvider";
import { AppIcon } from "./components/app-icon";
import { Header } from "./components/header";
import { ColorModeContextProvider } from "./contexts/color-mode";
import { UsersCreate, UsersEdit, UsersList, UsersShow } from "./pages/users";
import {
  CollectionsCreate,
  CollectionsEdit,
  CollectionsList,
  CollectionsShow,
} from "./pages/collections";
import { Login } from "./pages/login";
import {
  AccountBookFilled,
  AliwangwangFilled,
  FileImageFilled,
  FireFilled,
  PaperClipOutlined,
  TagFilled,
} from "@ant-design/icons";
import Title from "antd/es/typography/Title";
import {
  SBTTokensCreate,
  SBTTokensList,
  SBTTokensShow,
} from "./pages/tokens";
import { TonConnectUIProvider } from "@tonconnect/ui-react";
import { useOne, HttpError } from "@refinedev/core";
import {
  ActivitiesCreate,
  ActivitiesEdit,
  ActivitiesList,
  ActivitiesShow,
} from "./pages/activities";
import {
  PrototypeNftsCreate,
  PrototypeNftsEdit,
  PrototypeNftsList,
  PrototypeNftsShow,
} from "./pages/prototype-nfts";
import { RewardsList, RewardsShow } from "./pages/rewards";
import { RolesCreate, RolesEdit, RolesList, RolesShow } from "./pages/roles";
import axios, { AxiosInstance } from "axios";
import { SBTTokensEdit } from "./pages/tokens/edit";

export const API_URL = "https://tdp.tonbuilders.com";

const reactBoom = () => <Title level={4}>TDP Admin</Title>;

const axiosInstance = axios.create();

axiosInstance.interceptors.request.use((request: any) => {
  // Retrieve the token from local storage
  const token = localStorage.getItem("refine-auth");
  // Check if the header property exists
  if (request.headers) {
    // Set the Authorization header if it exists
    request.headers["Authorization"] = `Bearer ${token}`;
  } else {
    // Create the headers property if it does not exist
    request.headers = {
      Authorization: `Bearer ${token}`,
    };
  }

  return request;
});

type Permission = {
  id: string;
  route: string;
  method: string;
};

function App() {
  const { t, i18n } = useTranslation();

  const i18nProvider = {
    translate: (key: string, params: object) => t(key, params),
    changeLocale: (lang: string) => i18n.changeLanguage(lang),
    getLocale: () => i18n.language,
  };

  return (
    <BrowserRouter >
      <TonConnectUIProvider manifestUrl="https://tdp.tonbuilders.com/v1/manifest-ton-connect">
        <RefineKbarProvider>
          <ColorModeContextProvider>
            <Refine
          
              dataProvider={dataProvider(
                API_URL + "/v1/admin",
                axiosInstance as any
              )}
              notificationProvider={notificationProvider}
              authProvider={authProvider}
              i18nProvider={i18nProvider}
              accessControlProvider={{
                can: async ({ resource, action }) => {
                  const user = localStorage.getItem(USER_KEY);

                  if (!user) {
                    return { can: false };
                  }

                  const { permissions } = JSON.parse(user) as {
                    permissions: Permission[];
                  };

                  if (!permissions) {
                    return { can: false };
                  }

                  for (const permission of permissions) {
                    const route = permission.route.replace("/v1/admin/", "");
                    const routeParts = route.split("/");

                    // Map the HTTP method to an action
                    const methodToAction = {
                      GET: ["list", "show"],
                      POST: ["create"],
                      PUT: ["edit"],
                      PATCH: ["edit"],
                      DELETE: ["delete"],
                    } as { [key: string]: string[] };

                    // Check if the action matches the desired resource and action
                    if (
                      routeParts[0] === resource &&
                      methodToAction[permission.method].includes(action)
                    ) {
                      return { can: true };
                    }
                  }
                  return { can: false };
                },
              }}
              routerProvider={routerBindings}
              resources={[
                {
                  name: "users",
                  list: "/users",
                  create: "/users/create",
                  edit: "/users/edit/:id",
                  show: "/users/show/:id",
                  meta: {
                    canDelete: true,
                  },
                },
                {
                  name: "activities",

                  list: "/activities",
                  create: "/activities/create",
                  edit: "/activities/edit/:id",
                  show: "/activities/show/:id",
                  meta: {
                    canDelete: true,
                    icon: <FireFilled />,
                  },
                },
                {
                  name: "prototype-nfts",

                  list: "/prototype-nfts",
                  create: "/prototype-nfts/create",
                  edit: "/prototype-nfts/edit/:id",
                  show: "/prototype-nfts/show/:id",
                  meta: {
                    label: "NFT Templates",
                    canDelete: true,
                    icon: <TagFilled />,
                  },
                },
                {
                  name: "minted-nfts",
                  list: "/minted-nfts",
                  show: "/minted-nfts/show/:id",
                  edit: "/minted-nfts/edit/:id",

                  create: "/minted-nfts/create",
                  meta: {
                    canDelete: true,
                    icon: <AccountBookFilled />,
                  },
                },
                {
                  name: "collections",
                  list: "/collections",
                  create: "/collections/create",
                  edit: "/collections/edit/:id",
                  show: "/collections/show/:id",
                  meta: {
                    canDelete: true,
                    icon: <FileImageFilled />,
                  },
                },
                {
                  name: "rewards",
                  list: "/rewards",
                  show: "/rewards/show/:id",
                  meta: {
                    canDelete: true,
                    icon: <AliwangwangFilled />,
                  },
                },
                {
                  name: "roles",

                  list: "/roles",
                  create: "/roles/create",
                  edit: "/roles/edit/:id",
                  show: "/roles/show/:id",
                  meta: {
                    canDelete: true,
                    icon: <PaperClipOutlined />,
                  },
                },
              ]}
              options={{
                syncWithLocation: true,
                warnWhenUnsavedChanges: true,
              }}
            >
              <Routes>
                <Route
                  element={
                    <Authenticated fallback={<CatchAllNavigate to="/login" />}>
                      <ThemedLayoutV2
                        Header={() => <Header sticky />}
                        Sider={() => <ThemedSiderV2 Title={reactBoom} fixed />}
                        Title={({ collapsed }) => (
                          <ThemedTitleV2
                            collapsed={collapsed}
                            text="TDP Admin"
                            icon={<AppIcon />}
                          />
                        )}
                      >
                        <Outlet />
                      </ThemedLayoutV2>
                    </Authenticated>
                  }
                  path="/"
                >
                  <Route
                    index
                    element={<NavigateToResource resource="users" />}
                  />
                  <Route path="/users">
                    <Route index element={<UsersList />} />
                    <Route path="create" element={<UsersCreate />} />
                    <Route path="edit/:id" element={<UsersEdit />} />
                    <Route path="show/:id" element={<UsersShow />} />
                  </Route>

                  <Route path="/activities">
                    <Route index element={<ActivitiesList />} />
                    <Route path="create" element={<ActivitiesCreate />} />
                    <Route path="edit/:id" element={<ActivitiesEdit />} />
                    <Route path="show/:id" element={<ActivitiesShow />} />
                  </Route>

                  <Route path="/minted-nfts">
                    <Route index element={<SBTTokensList />} />
                    <Route path="create" element={<SBTTokensCreate />} />

                    <Route path="edit/:id" element={<SBTTokensEdit />} />

                    <Route path="show/:id" element={<SBTTokensShow />} />

                  </Route>

                  <Route path="/prototype-nfts">
                    <Route index element={<PrototypeNftsList />} />
                    <Route path="create" element={<PrototypeNftsCreate />} />
                    <Route path="edit/:id" element={<PrototypeNftsEdit />} />
                    <Route path="show/:id" element={<PrototypeNftsShow />} />
                  </Route>

                  <Route path="/rewards">
                    <Route index element={<RewardsList />} />
                    <Route path="show/:id" element={<RewardsShow />} />
                  </Route>

                  <Route path="/roles">
                    <Route index element={<RolesList />} />
                    <Route path="create" element={<RolesCreate />} />
                    <Route path="edit/:id" element={<RolesEdit />} />
                    <Route path="show/:id" element={<RolesShow />} />
                  </Route>

                  <Route path="/collections">
                    <Route index element={<CollectionsList />} />
                    <Route path="create" element={<CollectionsCreate />} />
                    <Route path="edit/:id" element={<CollectionsEdit />} />
                    <Route path="show/:id" element={<CollectionsShow />} />
                  </Route>

                  <Route path="*" element={<ErrorComponent />} />
                </Route>
                <Route
                  element={
                    <Authenticated fallback={<Outlet />}>
                      <NavigateToResource />
                    </Authenticated>
                  }
                >
                  <Route path="/login" element={<Login />} />
                </Route>
              </Routes>

              <RefineKbar />
              <UnsavedChangesNotifier />
            </Refine>
          </ColorModeContextProvider>
        </RefineKbarProvider>
      </TonConnectUIProvider>
    </BrowserRouter>
  );
}

export default App;
