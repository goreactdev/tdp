import React from "react";
import {
  IResourceComponentsProps,
  useApiUrl,
  useNavigation,
  useNotification,
} from "@refinedev/core";
import { useForm, getValueFromEvent } from "@refinedev/antd";
import { Button, Form, Input, Upload } from "antd";
import { useTonConnectUI } from "@tonconnect/ui-react";
import { Create } from "../../components/crud/create";

export const CollectionsCreate: React.FC<IResourceComponentsProps> = () => {
  const [tonConnectUI] = useTonConnectUI();

  const { goBack } = useNavigation();

  const { open } = useNotification();

  const [existing, setExisting] = React.useState(false);

  const { formProps, saveButtonProps } = useForm({
    resource: existing ? "existing-collection" : "collections",
    redirect: false,
    onMutationSuccess: async (data, _error, _context) => {

      if (existing) {
        goBack();
        return;
      }
      const transaction = {
        validUntil: Date.now() + 1000000,
        messages: [
          {
            address: data.data.contract_address,
            amount: "100000000",
            stateInit: data?.data.msg_body,
          },
        ],
      };

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


  const apiUrl = useApiUrl();

  return (
    <Create
      headerButtons={
        <Button onClick={() => setExisting(!existing)}>
          {existing ? "Create New" : "Add Existing"}
        </Button>
      }
      text={!existing ? 'Create New Collection' : "Add Existing Collection"}
      saveButtonProps={saveButtonProps}
    >
      <Form {...formProps} layout="vertical">
        {!existing && (
          <>
            <Form.Item
              label="Name"
              name={["display_name"]}
              rules={[
                {
                  required: true,
                },
              ]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              label="Description"
              name={["description"]}
              rules={[
                {
                  required: true,
                },
              ]}
            >
              <Input />
            </Form.Item>

            <Form.Item label="Cover Image">
              <Form.Item
                name="banner_image_url"
                valuePropName="fileList"
                getValueFromEvent={getValueFromEvent}
                noStyle
              >
                <Upload.Dragger
                  name="file"
                  headers={{
                    Authorization: `Bearer ${localStorage.getItem("refine-auth")}`,
                  }}
                  action={`${apiUrl}/media/upload`}
                      listType="picture"
                  maxCount={1}
                >
                  <p className="ant-upload-text">
                    Drag & drop a file in this area
                  </p>
                </Upload.Dragger>
              </Form.Item>
            </Form.Item>
            <Form.Item label="Collection Image">
              <Form.Item
                name="collection_image_url"
                valuePropName="fileList"
                getValueFromEvent={getValueFromEvent}
                noStyle
                rules={[
                  {
                    required: true,
                  },
                ]}
              >
                <Upload.Dragger
                  name="file"
                  headers={{
                    Authorization: `Bearer ${localStorage.getItem("refine-auth")}`,
                  }}
                  action={`${apiUrl}/media/upload`}
                      listType="picture"
                  maxCount={1}
                >
                  <p className="ant-upload-text">
                    Drag & drop a file in this area
                  </p>
                </Upload.Dragger>
              </Form.Item>
            </Form.Item>
          </>
        )}

        {existing && (
          <Form.Item
          label="Friendly Address"
          name={["friendly_address"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input />
        </Form.Item>
        )}        

        <Form.Item
          label="Default rating points"
          name={["default_weight"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input />
        </Form.Item>
      </Form>
    </Create>
  );
};
