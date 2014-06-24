/**
 * @jsx React.DOM
 */

window.app = (function(window, document, undefined){
  var Page = React.createClass({displayName: 'Page',
    render: function() {
      return  React.DOM.div( {id:"page"}, 
                Title(null ),
                Grid(null ),
                Footer(null )
              );
    }
  });

  var Title = React.createClass({displayName: 'Title',
    render: function() {
      return React.DOM.h1(null, "Tarcísio Gruppi");
    }
  });

  var Grid = React.createClass({displayName: 'Grid',
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
        return  React.DOM.div( {id:"grid-loader", className:"csspinner double-up"});
      } else if (this.state.items === false || this.state.items.length === 0) {
        return  React.DOM.div( {id:"grid-error"}, 
                  React.DOM.p(null, "Sadly it was not possible to load this page. You can try again later."),
                  React.DOM.p(null, "Feel free to contact me at ", React.DOM.a( {href:"mainto:txgruppi@gmail.com"}, "txgruppi@gmail.com"),".")
                );
      } else {
        var items = this.state.items.map(function(link){
          return GridItem( {url:link.url, title:link.title, image:link.image} )
        });
        return React.DOM.ul( {id:"grid"}, items);
      }
    }
  });

  var GridItem = React.createClass({displayName: 'GridItem',
    propTypes: {
      url: React.PropTypes.string.isRequired,
      title: React.PropTypes.string.isRequired,
      image: React.PropTypes.string.isRequired
    },

    render: function() {
      return  React.DOM.li(null, 
                React.DOM.a( {href:this.props.url, title:this.props.title}, 
                  React.DOM.img( {src:this.props.image, alt:this.props.title} )
                )
              );
    }
  });

  var Footer = React.createClass({displayName: 'Footer',
    render: function() {
      return  React.DOM.p(null, 
                "Served by ", React.DOM.a( {href:"http://martini.codegangsta.io/", title:"martini"}, "martini"),React.DOM.br(null),
                "(some) Icons by ", React.DOM.a( {href:"http://simpleicons.org/"}, "Simple Icons"),React.DOM.br(null),
                "Background by ", React.DOM.a( {href:"http://subtlepatterns.com/"}, "Subtle Patterns")
              );
    }
  });

  window.pageView = React.renderComponent(Page(null ), document.getElementById("stage"));
})(window, document);
