/**
 * @jsx React.DOM
 */

window.app = (function(window, document, undefined){
  var Page = React.createClass({
    render: function() {
      return  <div id="page">
                <Title />
                <Grid />
                <Footer />
              </div>;
    }
  });

  var Title = React.createClass({
    render: function() {
      return <h1>Tarcísio Gruppi</h1>;
    }
  });

  var Grid = React.createClass({
    getInitialState: function(){
      return {
        items: null
      };
    },

    componentDidMount: function() {
      var that = this;
      nwt.io("/api/links").success(function(event){
        if (event.obj instanceof Array) {
          that.setState({
            items: event.obj
          });
        } else {
          that.setState({
            items: false
          })
        }
      }).get();
    },

    render: function() {
      if (this.state.items === null) {
        return  <div id="grid-loader" className="csspinner double-up"></div>;
      } else if (this.state.items === false || this.state.items.length === 0) {
        return  <div id="grid-error">
                  <p>Sadly it was not possible to load this page. You can try again later.</p>
                  <p>Feel free to contact me at <a href="mainto:txgruppi@gmail.com">txgruppi@gmail.com</a>.</p>
                </div>;
      } else {
        var items = this.state.items.map(function(link){
          return <GridItem url={link.url} title={link.title} image={link.image} />
        });
        return <ul id="grid">{items}</ul>;
      }
    }
  });

  var GridItem = React.createClass({
    propTypes: {
      url: React.PropTypes.string.isRequired,
      title: React.PropTypes.string.isRequired,
      image: React.PropTypes.string.isRequired
    },

    render: function() {
      return  <li>
                <a href={this.props.url} title={this.props.title}>
                  <img src={this.props.image} alt={this.props.title} />
                </a>
              </li>;
    }
  });

  var Footer = React.createClass({
    render: function() {
      return  <p>
                Served by <a href="http://martini.codegangsta.io/" title="martini">martini</a><br/>
                (some) Icons by <a href="http://simpleicons.org/">Simple Icons</a><br/>
                Background by <a href="http://subtlepatterns.com/">Subtle Patterns</a>
              </p>;
    }
  });

  window.pageView = React.renderComponent(<Page />, document.getElementById("stage"));
})(window, document);
