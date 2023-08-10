import React from "react";
import {
  LoginPageProps,
  LoginFormTypes,
  useLink,
  useRouterType,
  useActiveAuthProvider,
  useLogin,
  useTranslate,
  useRouterContext,
} from "@refinedev/core";
import { ThemedTitle } from "@refinedev/antd";
import {
  bodyStyles,
  containerStyles,
  headStyles,
  layoutStyles,
  titleStyles,
} from "./styles";
import {
  Row,
  Col,
  Layout,
  Card,
  Typography,
  Form,
  Input,
  Button,
  Checkbox,
  CardProps,
  LayoutProps,
  Divider,
  FormProps,
  theme,
} from "antd";

import { TonConnectButton, TonProofItemReplySuccess, useTonConnectUI } from "@tonconnect/ui-react";
import { axiosInstance } from "@refinedev/simple-rest";
import { ProofCheckPayload } from "../../../../authProvider";
import { API_URL } from "../../../../App";

const { Text, Title } = Typography;
const { useToken } = theme;

type LoginProps = LoginPageProps<LayoutProps, CardProps, FormProps>;
/**
 * **refine** has a default login page form which is served on `/login` route when the `authProvider` configuration is provided.
 *
 * @see {@link https://refine.dev/docs/ui-frameworks/antd/components/antd-auth-page/#login} for more details.
 * 
 * 
 */


export interface ProofPayloadResponse {
  payload: string
}


export const LoginPage: React.FC<LoginProps> = ({
  providers,
  registerLink,
  forgotPasswordLink,
  rememberMe,
  contentProps,
  wrapperProps,
  renderContent,
  formProps,
  title,
}) => {
  const { token } = useToken();
  const [form] = Form.useForm<LoginFormTypes>();
  const translate = useTranslate();
  const routerType = useRouterType();
  const Link = useLink();
  const { Link: LegacyLink } = useRouterContext();

  const { mutate } = useLogin<ProofCheckPayload>({
    v3LegacyAuthProviderCompatible: false
  });
  

  const [tonConnectUI] = useTonConnectUI()
  const [payloadData, setPayloadData] = React.useState<ProofPayloadResponse | null>(null)

  React.useEffect(() => {
    tonConnectUI.disconnect()
  }, [])

  React.useEffect(() => {
    const proofPayloadQuery = async () => {
      const { data } = await axiosInstance.get(API_URL + '/v1/ton-connect/generate-payload')
      setPayloadData(data)
    }
    proofPayloadQuery()
  }, [])


  if (!payloadData) {
    tonConnectUI.setConnectRequestParameters(null)
  } else {
    tonConnectUI.setConnectRequestParameters({
      state: 'ready',
      value: { tonProof: payloadData.payload  },
    })
  }

  React.useEffect(() => {
    tonConnectUI.onStatusChange(async (wallet) => {
      if (
        wallet?.connectItems?.tonProof &&
        'proof' in wallet.connectItems.tonProof
      ) {

         mutate({
          address: wallet.account.address,
          network: wallet.account.chain,
          proof: wallet.connectItems.tonProof.proof,
        })
      }
    })
  }, [])

  const PageTitle =
    title === false ? null : (
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          marginBottom: "32px",
          fontSize: "20px",
        }}
      >
        {title ?? <ThemedTitle collapsed={false} />}
      </div>
    );

  const CardTitle = (
    <Title
      level={3}
      style={{
        color: token.colorPrimaryTextHover,
        ...titleStyles,
      }}
    >
      {translate("pages.login.title", "Sign in to your account")}
    </Title>
  );
  const CardContent = (
    <Card
      title={CardTitle}
      headStyle={headStyles}
      bodyStyle={bodyStyles}
      style={{
        ...containerStyles,
        backgroundColor: token.colorBgElevated,
      }}
      {...(contentProps ?? {})}
    >
      <div style={{display: "flex", justifyContent: "center", width: "100%"}}>
        <TonConnectButton  />
      </div>
    </Card>
  );

  return (
    <Layout style={layoutStyles} {...(wrapperProps ?? {})}>
      <Row
        justify="center"
        align="middle"
        style={{
          height: "100vh",
        }}
      >
        <Col xs={22}>
          {renderContent ? (
            renderContent(CardContent, PageTitle)
          ) : (
            <>
              {PageTitle}
              {CardContent}
            </>
          )}
        </Col>
      </Row>
    </Layout>
  );
};