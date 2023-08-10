import React from "react";
import { AuthPageProps } from "@refinedev/core";
import {
  LoginPage,
} from "./components";
import { CardProps, FormProps, LayoutProps } from "antd";

export type AuthProps = AuthPageProps<LayoutProps, CardProps, FormProps> & {
  renderContent?: (
    content: React.ReactNode,
    title: React.ReactNode
  ) => React.ReactNode;
  title?: React.ReactNode;
};

export const AuthPage: React.FC<AuthProps> = (props) => {
  return <LoginPage {...props} />;
};
