import {
  QueryClient,
  QueryClientProvider,
  useQuery,
} from "@tanstack/react-query";
import { Theme, Spinner, Callout, Table, Text } from "@radix-ui/themes";
import { InfoCircledIcon } from "@radix-ui/react-icons";

import "@radix-ui/themes/styles.css";
import "./app.css";
import { Line, LineChart, XAxis, YAxis } from "recharts";

// Create a client
const queryClient = new QueryClient();

function App() {
  return (
    // Provide the client to your App
    <QueryClientProvider client={queryClient}>
      <Theme>
        <Page>
          <Text size="7">Table</Text>
          <Fingerprints />
          <Text size="7">Graph</Text>
          <FingerprintGraphs />
        </Page>
      </Theme>
    </QueryClientProvider>
  );
}

const BASE_URL = "http://localhost:8080/api";

const Page = ({ children }: { children: React.ReactNode }) => (
  <div className="page">
    <div className="header">
      <Text size="9">Fingerprints</Text>
    </div>
    {children}
  </div>
);

// Data is of the format:
// {
//   "fingerprints": [{ input: string, timestamp: string }]
//   "current_page": int,
//   "total_pages": int,
// }
const fetchFingerprints = async () => {
  const response = await fetch(`${BASE_URL}/fingerprints`);
  return response.json();
};

const Fingerprints = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["fingerprints"],
    queryFn: fetchFingerprints,
  });

  if (isLoading) return <Spinner />;
  if (error)
    return (
      <Callout.Root>
        <Callout.Icon>
          <InfoCircledIcon />
        </Callout.Icon>
        <Callout.Text>
          An error occurred while fetching fingerprints. Please try again.
        </Callout.Text>
      </Callout.Root>
    );

  return (
    <Table.Root>
      <Table.Header>
        <Table.Row>
          <Table.ColumnHeaderCell>Input</Table.ColumnHeaderCell>
          <Table.ColumnHeaderCell>Timestamp</Table.ColumnHeaderCell>
        </Table.Row>
      </Table.Header>
      <Table.Body>
        {data.fingerprints.map(
          (fingerprint: { input: string; timestamp: string }, i) => (
            <Table.Row key={i}>
              <Table.Cell>{fingerprint.input}</Table.Cell>
              <Table.Cell>{fingerprint.timestamp}</Table.Cell>
            </Table.Row>
          ),
        )}
      </Table.Body>
    </Table.Root>
  );
};

// Data is of the format:
// [{ timestamp: string, count: int }]
const fetchFingerprintCount = async () => {
  const response = await fetch(`${BASE_URL}/fingerprints/count`);
  return response.json();
};

const FingerprintGraphs = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["fingerprintCount"],
    queryFn: fetchFingerprintCount,
  });

  if (isLoading) return <Spinner />;
  if (error)
    return (
      <Callout.Root>
        <Callout.Icon>
          <InfoCircledIcon />
        </Callout.Icon>
        <Callout.Text>
          An error occurred while fetching fingerprint count. Please try again.
        </Callout.Text>
      </Callout.Root>
    );

  console.log(data);

  // Render a rechart line graph.
  return (
    <div>
      <LineChart width={600} height={300} data={data}>
        <XAxis dataKey="timestamp" />
        <YAxis />
        <Line type="monotone" dataKey="count" stroke="#8884d8" />
      </LineChart>
    </div>
  );
};

export default App;
