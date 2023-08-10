import React from "react";
import {
  useNavigation,
  useTranslate,
  userFriendlyResourceName,
  useRefineContext,
  useRouterType,
  useResource,
  useBack,
} from "@refinedev/core";
import {
  Breadcrumb,
  SaveButton,
  PageHeader,
  SaveButtonProps,
  useDrawerForm,
  useForm,
} from "@refinedev/antd";
import { Card, Space, Spin } from "antd";
import type { CreateProps } from "@refinedev/antd";
export const Create: React.FC<CreateProps & {text?:string}> = ({
  saveButtonProps: saveButtonPropsFromProps,
  children,
  resource: resourceFromProps,
  isLoading = false,
  breadcrumb: breadcrumbFromProps,
  wrapperProps,
  headerProps,
  contentProps,
  headerButtonProps,
  headerButtons,
  footerButtonProps,
  footerButtons,
  goBack: goBackFromProps,
  text = 'Mint SBT',
}) => {
  const { options: { breadcrumb: globalBreadcrumb } = {} } = useRefineContext();


  const routerType = useRouterType();
  const back = useBack();
  const { goBack } = useNavigation();


 
  const { resource, action } = useResource(resourceFromProps);

  const breadcrumb =
    typeof breadcrumbFromProps === "undefined"
      ? globalBreadcrumb
      : breadcrumbFromProps;

  const saveButtonProps: SaveButtonProps = {
    ...(isLoading ? { disabled: true } : {}),
    ...saveButtonPropsFromProps,
    htmlType: "submit",
  };


  const defaultFooterButtons = (
    <>      
      <SaveButton {...saveButtonProps} />
      
    </>
  );

  return (
    <div {...(wrapperProps ?? {})}>
      <PageHeader
        ghost={false}
        backIcon={goBackFromProps}
        onBack={
          action !== "list" || typeof action !== "undefined"
            ? routerType === "legacy"
              ? goBack
              : back
            : undefined
        }
        title={text}
        breadcrumb={
          typeof breadcrumb !== "undefined" ? (
            <>{breadcrumb}</> ?? undefined
          ) : (
            <Breadcrumb />
          )
        }
        extra={
          <Space wrap {...(headerButtonProps ?? {})}>
            {headerButtons
              ? typeof headerButtons === "function"
                ? headerButtons({
                    defaultButtons: null,
                  })
                : headerButtons
              : null}
          </Space>
        }
        {...(headerProps ?? {})}
      >
        <Spin spinning={isLoading}>
          <Card
            bordered={false}
            actions={[
              <Space
                key="action-buttons"
                style={{ float: "right", marginRight: 24 }}
                {...(footerButtonProps ?? {})}
              >
                {footerButtons
                  ? typeof footerButtons === "function"
                    ? footerButtons({
                        defaultButtons: defaultFooterButtons,
                        saveButtonProps: saveButtonProps,
                      })
                    : footerButtons
                  : defaultFooterButtons}
              </Space>,
            ]}
            {...(contentProps ?? {})}
          >
            {children}
          </Card>
        </Spin>
      </PageHeader>
    </div>
  );
};
