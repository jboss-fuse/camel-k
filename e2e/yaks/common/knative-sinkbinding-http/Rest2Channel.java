import org.apache.camel.builder.RouteBuilder;

public class Rest2Channel extends RouteBuilder {
    public void configure() throws Exception {
        rest("/")
            .put("/foo/new")
            .to("knative:channel/messages");
    }
}