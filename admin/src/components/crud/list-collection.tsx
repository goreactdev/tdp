import React from "react";
import {
  useTranslate,
  userFriendlyResourceName,
  useRefineContext,
  useRouterType,
  useResource,
} from "@refinedev/core";
import {
  Breadcrumb,
  CreateButton,
  CreateButtonProps,
  PageHeader,
} from "@refinedev/antd";
import { Space } from "antd";
import type { ListProps } from "@refinedev/antd";

export const List: React.FC<ListProps> = ({
  canCreate,
  title,
  children,
  createButtonProps: createButtonPropsFromProps,
  resource: resourceFromProps,
  wrapperProps,
  contentProps,
  headerProps,
  breadcrumb: breadcrumbFromProps,
  headerButtonProps,
  headerButtons,
}) => {
  const translate = useTranslate();
  const { options: { breadcrumb: globalBreadcrumb } = {} } = useRefineContext();

  const routerType = useRouterType();

  const { resource } = useResource(resourceFromProps);

  const isCreateButtonVisible =
    canCreate ??
    ((resource?.canCreate ?? !!resource?.create) || createButtonPropsFromProps);

  const breadcrumb =
    typeof breadcrumbFromProps === "undefined"
      ? globalBreadcrumb
      : breadcrumbFromProps;

  const createButtonProps: CreateButtonProps | undefined = isCreateButtonVisible
    ? {
        size: "middle",
        resource:
          routerType === "legacy"
            ? resource?.route
            : resource?.identifier ?? resource?.name,
        ...createButtonPropsFromProps,
      }
    : undefined;

  const defaultExtra = isCreateButtonVisible ? (
    <CreateButton {...createButtonProps} />
  ) : null;

  return (
    <div {...(wrapperProps ?? {})}>
      <PageHeader
        ghost={false}
        title={
          title ??
          translate(
            `${resource?.name}.titles.list`,
            userFriendlyResourceName(
              resource?.meta?.label ??
                resource?.options?.label ??
                resource?.label ??
                resource?.name,
              "plural"
            )
          )
        }
        extra={
          headerButtons ? (
            <Space wrap {...headerButtonProps}>
              {typeof headerButtons === "function"
                ? headerButtons({
                    defaultButtons: defaultExtra,
                    createButtonProps,
                  })
                : headerButtons}
            </Space>
          ) : (
            defaultExtra
          )
        }
        breadcrumb={
          typeof breadcrumb !== "undefined" ? (
            <>{breadcrumb}</> ?? undefined
          ) : (
            <Breadcrumb />
          )
        }
        {...(headerProps ?? {})}
      >
        <div {...(contentProps ?? {})}>{children}</div>
      </PageHeader>
    </div>
  );
};
