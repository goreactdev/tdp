import { DownOutlined, SearchOutlined } from "@ant-design/icons";
import { RefineThemedLayoutV2HeaderProps } from "@refinedev/antd";
import {
  useGetIdentity,
  useGetLocale,
  useList,
  useSetLocale,
} from "@refinedev/core";
import {
  Avatar,
  Button,
  Dropdown,
  Layout as AntdLayout,
  MenuProps,
  Space,
  Switch,
  theme,
  Typography,
  AutoComplete,
  Input,
} from "antd";
import React, { useContext, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { ColorModeContext } from "../../contexts/color-mode";
import { Link } from "react-router-dom";

const { Text } = Typography;
const { useToken } = theme;

type IUser = {
  id: number;
  name: string;
  avatar: string;
};

const renderTitle = (title: string) => {
  return (
    <Text strong style={{ fontSize: "16px" }}>
      {title}
    </Text>
  );
};

// To be able to customize the option item
const renderItem = (
  title: string,
  friendlyAddr: string,
  resource: string,
  id: number
) => {
  return {
    value: title,
    label: (
      <Link to={`/${resource}/show/${id}`}>
        <Space.Compact direction="vertical">
          <Text strong>{title}</Text>
          <Text>{friendlyAddr}</Text>
        </Space.Compact>
      </Link>
    ),
  };
};

export const Header: React.FC<RefineThemedLayoutV2HeaderProps> = ({
  sticky,
}) => {
  const [value, setValue] = useState<string>("");
  const [options, setOptions] = useState<any>([]);

  const { refetch: refetchItems } = useList({
    resource: "minted-nfts",
    filters: [{ field: "name", operator: "contains", value }],
    queryOptions: {
      enabled: false,
      onSuccess: (data) => {
        const nftsOptionGroup = data.data.map((item) =>
          renderItem(
            item.name,
            item.friendly_address,
            "minted-nfts",
            item.id as any
          )
        );
        if (nftsOptionGroup.length > 0) {
          setOptions([
            {
              label: renderTitle("NFTs"),
              options: nftsOptionGroup,
            },
          ]);
        }
      },
    },
  });

  useEffect(() => {
    refetchItems();
    setOptions([]);
  }, [value]);

  const { token } = useToken();
  const { i18n } = useTranslation();
  const locale = useGetLocale();
  const changeLanguage = useSetLocale();
  const { data: user } = useGetIdentity<IUser>();
  const { mode, setMode } = useContext(ColorModeContext);

  const currentLocale = locale();

  const menuItems: MenuProps["items"] = [...(i18n.languages || [])]
    .sort()
    .map((lang: string) => ({
      key: lang,
      onClick: () => changeLanguage(lang),
      icon: (
        <span style={{ marginRight: 8 }}>
          <Avatar size={16} src={`/images/flags/${lang}.svg`} />
        </span>
      ),
      label: lang === "en" ? "English" : "German",
    }));

  const headerStyles: React.CSSProperties = {
    backgroundColor: token.colorBgElevated,
    display: "flex",
    justifyContent: "flex-end",
    alignItems: "center",
    padding: "0px 24px",
    height: "64px",
  };

  if (sticky) {
    headerStyles.position = "sticky";
    headerStyles.top = 0;
    headerStyles.zIndex = 1;
  }

  const isMobile = window.innerWidth < 768;

  return (
    <AntdLayout.Header style={headerStyles}>
      <AutoComplete
        style={{ width: "100%", maxWidth: "550px" }}
        filterOption={false}
        options={options}
        onSearch={(value: string) => setValue(value)}
      >
        <Input
          size="large"
          placeholder="Search prototypes..."
          suffix={<SearchOutlined />}
        />
      </AutoComplete>
      <Space>
        <Dropdown
          menu={{
            items: menuItems,
            selectedKeys: currentLocale ? [currentLocale] : [],
          }}
        >
          <Button type="text">
            <Space>
              <Avatar size={16} src={`/images/flags/${currentLocale}.svg`} />
              {currentLocale === "en" ? "English" : "German"}
              <DownOutlined />
            </Space>
          </Button>
        </Dropdown>
        <Switch
          checkedChildren="🌛"
          unCheckedChildren="🔆"
          onChange={() => setMode(mode === "light" ? "dark" : "light")}
          defaultChecked={mode === "dark"}
        />
        <Space style={{ marginLeft: "8px" }} size="middle">
          {!isMobile && user?.name && <Text strong>{user.name}</Text>}
          {user?.avatar && <Avatar src={user?.avatar} alt={user?.name} />}
        </Space>
      </Space>
    </AntdLayout.Header>
  );
};
