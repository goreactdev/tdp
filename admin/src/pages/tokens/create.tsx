import React from "react";
import {
  IResourceComponentsProps,
  useList,
  useNavigation,
} from "@refinedev/core";
import {
  ImportButton,
  getValueFromEvent,
  useForm,
  useImport,
} from "@refinedev/antd";
import {
  Form,
  Input,
  DatePicker,
  AutoComplete,
  Space,
  Typography,
  Upload,
  Button,
} from "antd";
import dayjs from "dayjs";
import { Create } from "../../components/crud/create";
import { Link } from "react-router-dom";
import { useState } from "react";
import { useEffect } from "react";
import {
  InboxOutlined,
  SearchOutlined,
  UploadOutlined,
} from "@ant-design/icons";
import { useTonConnectUI } from "@tonconnect/ui-react";
import {
  useOne,
  useApiUrl,
  useCustomMutation,
  useNotification,
} from "@refinedev/core";
import Dragger from "antd/es/upload/Dragger";
import { file2Base64 } from "@refinedev/core";

const { Text } = Typography;

// To be abl
// To be able to customize the o
export const SBTTokensCreate: React.FC<IResourceComponentsProps> = () => {
  const [tonConnectUI] = useTonConnectUI();

  const { goBack } = useNavigation();

  const { open } = useNotification();

  const { formProps, saveButtonProps } = useForm({
    redirect: false,
    onMutationError: (error, _context) => {
      open?.({
        message: "Error",
        description: error.message,
        type: "error",
      });
    },
    onMutationSuccess: async (data, _error, _context) => {
      let transaction;

      if (data.data.fee_for_tx) {
        transaction = {
          validUntil: Date.now() + 1000000,
          messages: [
            {
              address: data.data.collection_address,
              amount: String(parseFloat(data.data.fee_for_tx) * 1e9),
              payload: data.data.msg_body,
            },
          ],
        };
      } else {
        transaction = {
          validUntil: Date.now() + 1000000,
          messages: [
            {
              address: data.data.collection_address,
              amount: "20000000",
              payload: data.data.msg_body,
            },
          ],
        };
      }

      const result = await tonConnectUI.sendTransaction(transaction);

      if (result) {
        open?.({
          message: "Transaction sent",
          description: "Transaction sent successfully",
          type: "success",
        });

        goBack();

        return;
      }

      open?.({
        message: "Transaction failed",
        description: "Transaction failed",
        type: "error",
      });
    },
  });

  const [options, setOptions] = useState<any>([]);
  const [tokenOptions, setTokenOptions] = useState<any>([]);

  const [searchData, setSearchData] = useState<{
    user_address: string;
    collection_address: string;
    activity_name: string;
    manifest: string;
  }>({
    user_address: "",
    collection_address: "",
    activity_name: "",
    manifest: "",
  });

  const apiUrl = useApiUrl();

  const { refetch: refetchItems } = useList({
    resource: "users",
    filters: [
      { field: "name", operator: "contains", value: searchData.user_address },
    ],
    queryOptions: {
      enabled: false,
      onSuccess: (data) => {
        const usersOptionGroups = data.data.map((item) => ({
          value: item.username,
          label: (
            <Space.Compact direction="vertical">
              <Text strong>
                {item.first_name} {item.last_name}
              </Text>
              <Text>{item.username}</Text>
            </Space.Compact>
          ),
        }));
        if (usersOptionGroups.length > 0) {
          setOptions([
            {
              label: (
                <Text strong style={{ fontSize: "16px" }}>
                  Users
                </Text>
              ),
              options: usersOptionGroups,
            },
          ]);
        }
      },
    },
  });

  const { refetch: refetchTokens } = useList({
    resource: "collections",

    filters: [
      {
        field: "name",
        operator: "contains",
        value: searchData.collection_address,
      },
      {
        field:"owner_address",
        operator: 'null',
        value: "true",
      }
    ],
    queryOptions: {
      enabled: false,
      onSuccess: (data) => {
        const collectionsOptionGroup = data.data.map((item) => ({
          value: item.friendly_address,
          label: (
            <Space.Compact direction="vertical">
              <Text strong>{item.name}</Text>
              <Text>{item.friendly_address}</Text>
            </Space.Compact>
          ),
        }));
        if (collectionsOptionGroup.length > 0) {
          setTokenOptions([
            {
              label: (
                <Text strong style={{ fontSize: "16px" }}>
                  Collections
                </Text>
              ),
              options: collectionsOptionGroup,
            },
          ]);
        }
      },
    },
  });

  const [manifest, setManifest] = useState<any>([]);

  const { refetch: refetchManifest } = useList({
    resource: "prototype-nfts",
    filters: [
      {
        field: "name",
        operator: "contains",
        value: searchData.manifest,
      },
    ],
    queryOptions: {
      enabled: false,
      onSuccess: (data) => {
        const tokensOptionGroup = data.data.map((item) => ({
          value: item.id + ": " + item.display_name,
          label: (
            <Space.Compact direction="vertical">
              <Text strong>{item.id + ": " +  item.display_name}</Text>
            </Space.Compact>
          ),
        }));
        if (tokensOptionGroup.length > 0) {
          setManifest([
            {
              label: (
                <Text strong style={{ fontSize: "16px" }}>
                  Tokens
                </Text>
              ),

              options: tokensOptionGroup,
            },
          ]);
        }
      },
    },
  });

  useEffect(() => {
    refetchManifest();
    setManifest([]);
  }, [searchData.manifest]);

  useEffect(() => {
    refetchTokens();
    setTokenOptions([]);
  }, [searchData.collection_address]);

  useEffect(() => {
    refetchItems();
    setOptions([]);
  }, [searchData.user_address]);

  const [isMassMint, setIsMassMint] = useState(false);

  return (
    <Create
      headerButtons={() => (
        <Button onClick={() => setIsMassMint(!isMassMint)}>
          {isMassMint ? "Single mint" : "Mass mint"}
        </Button>
      )}
      saveButtonProps={saveButtonProps}
    >
      <Form {...formProps} layout="vertical">
        {!isMassMint && (
          <Form.Item
            label="User"
            name={["user_address"]}
            rules={[
              {
                required: true,
              },
            ]}
          >
            <AutoComplete
              style={{ width: "100%", maxWidth: "550px" }}
              filterOption={false}
              options={options}
              value={searchData.user_address}
              onChange={(value: string) =>
                setSearchData({ ...searchData, user_address: value })
              }
              onSearch={(value: string) =>
                setSearchData({ ...searchData, user_address: value })
              }
            >
              <Input
                size="large"
                value={searchData.user_address}
                placeholder="Search users"
                suffix={<SearchOutlined />}
              />
            </AutoComplete>
          </Form.Item>
        )}

        {isMassMint && (
          <Form.Item label="CSV file">
            <Form.Item
              name={["csv_file"]}
              valuePropName="fileList"
              getValueFromEvent={getValueFromEvent}
              noStyle
            >
              <Upload.Dragger
                name="file"
                action={`${apiUrl}/csv/upload`}
                listType="text"
                maxCount={1}
              >
                <p className="ant-upload-text">
                  Drag & drop a file in this area <br />
                  <Text type="warning"> Max is 80 addresses per upload</Text>
                </p>
              </Upload.Dragger>
            </Form.Item>
          </Form.Item>
        )}

        <Form.Item
          label="Manifest NFT"
          name={["meta_json_id"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <AutoComplete
            style={{ width: "100%", maxWidth: "550px" }}
            filterOption={false}
            options={manifest}
            value={searchData.manifest}
            onChange={(value: string) => {
              setSearchData({ ...searchData, manifest: value });
            }}
            onSearch={(value: string) =>
              setSearchData({ ...searchData, manifest: value })
            }
          >
            <Input
              size="large"
              placeholder="Search manifest"
              suffix={<SearchOutlined />}
            />
          </AutoComplete>
        </Form.Item>
        <Form.Item
          label="SBT Collection"
          name={["collection_address"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <AutoComplete
            style={{ width: "100%", maxWidth: "550px" }}
            filterOption={false}
            options={tokenOptions}
            value={searchData.collection_address}
            onChange={(value: string) => {
              setSearchData({ ...searchData, collection_address: value });
            }}
            onSearch={(value: string) =>
              setSearchData({ ...searchData, collection_address: value })
            }
          >
            <Input
              size="large"
              value={searchData.collection_address}
              placeholder="Search collections"
              suffix={<SearchOutlined />}
            />
          </AutoComplete>
        </Form.Item>
      </Form>
    </Create>
  );
};
